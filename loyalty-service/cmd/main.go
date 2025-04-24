package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/loyalty-service/internal/kafka"
	"github.com/paranoiachains/loyalty-api/pkg/flags"
)

func main() {
	kafka.StartKafkaServices(kafka.Input, kafka.Processed)

	r := gin.New()
	r.Run(flags.AccrualSystemAddress)
}
