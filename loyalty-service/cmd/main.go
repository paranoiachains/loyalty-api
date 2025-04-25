package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/paranoiachains/loyalty-api/pkg/flags"
	"github.com/paranoiachains/loyalty-api/pkg/messaging"
)

func main() {
	// init kafka services
	messaging.LoyaltyKafka = messaging.InitLoyaltyKafka()
	messaging.LoyaltyKafka.Start(context.Background())
	for v := range messaging.LoyaltyKafka.Receive() {
		fmt.Printf("%v\n", v)
	}

	r := gin.New()
	r.Run(flags.AccrualSystemAddress)
}
