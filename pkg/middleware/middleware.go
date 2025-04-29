package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/paranoiachains/loyalty-api/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		encoding := c.Request.Header.Get("Accept-Encoding")
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Log.Error("read body", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		c.Next()

		duration := time.Since(start)

		logger.Log.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Duration("duration", duration),
			zap.String("Accept-Encoding", encoding),
			zap.String("Request Body", string(body)),
		)
		logger.Log.Info("HTTP Response",
			zap.Int("status", c.Writer.Status()),
			zap.String("Content-Type", c.Writer.Header().Get("Content-Type")),
		)
	}
}

func shouldCompress(req *http.Request) bool {
	return strings.Contains(req.Header.Get("Accept-Encoding"), "gzip")
}

func shouldDecompress(req *http.Request) bool {
	return strings.Contains(req.Header.Get("Content-Encoding"), "gzip")
}

func Compression() gin.HandlerFunc {
	return func(c *gin.Context) {
		// decompress request
		if shouldDecompress(c.Request) {
			r, err := gzip.NewReader(c.Request.Body)
			if err == nil {
				defer r.Close()
				body, _ := io.ReadAll(r)
				c.Request.Body = io.NopCloser(bytes.NewReader(body))
				c.Request.Header.Del("Content-Encoding")
				logger.Log.Info("decompressed request", zap.Bool("ok", true))
			}
		}

		// if client doesn't accept gzip, skip
		if !shouldCompress(c.Request) {
			c.Next()
			return
		}

		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()
		c.Writer = &gzipWriter{ResponseWriter: c.Writer, writer: gz}

		c.Header("Content-Encoding", "gzip")
		c.Next()
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := &struct {
			jwt.RegisteredClaims
			UserID int64 `json:"user_id"`
		}{}
		_, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret_key"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}

var (
	limiters = make(map[int64]*rate.Limiter)
	mu       sync.Mutex
)

// возвращает лимитер для пользователя
func getLimiter(userID int64) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[userID]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(1*time.Second), 10)
		limiters[userID] = limiter
	}
	return limiter
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("userID")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID := val.(int64)

		limiter := getLimiter(userID)
		if !limiter.Allow() {
			c.Header("Retry-After", "60")
			c.String(http.StatusTooManyRequests, "No more than 10 requests per minute allowed")
			return
		}

		c.Next()
	}
}
