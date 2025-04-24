package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/pkg/flags"
)

func main() {
	r := gin.New()
	r.Run(flags.AccrualSystemAddress)
}
