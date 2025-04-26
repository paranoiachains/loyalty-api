package flags

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
)

var (
	RunAddress           string
	DatabaseDSN          string
	LoyaltyDatabaseDSN   string
	AccrualSystemAddress string
)

type Environment struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	OrderDatabaseDSN     string `env:"DB_DSN"`
	LoyaltyDatabaseDSN   string `env:"ACCRUAL_DB_DSN"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func init() {
	// service's own flag set
	accruals := flag.NewFlagSet("", flag.ExitOnError)

	// variables that were parsed from environment
	var parsedEnv Environment

	// parse variables only once
	var once sync.Once
	once.Do(func() {
		accruals.StringVar(&RunAddress, "a", ":8080", "service address and port")
		accruals.StringVar(&DatabaseDSN, "d", "postgresql://postgres:postgres@postgres/order_service?sslmode=disable", "database connection uri")
		accruals.StringVar(&LoyaltyDatabaseDSN, "dl", "postgresql://postgres:postgres@postgres/loyalty_service?sslmode=disable", "db connection uri")
		accruals.StringVar(&AccrualSystemAddress, "r", ":8080", "accrual system address")
		accruals.Parse(os.Args[1:])

		err := env.Parse(&parsedEnv)
		if err != nil {
			log.Fatal(err)
		}

		// prioritize environment variables over flags
		if parsedEnv.RunAddress != "" {
			RunAddress = parsedEnv.RunAddress
		}
		if parsedEnv.OrderDatabaseDSN != "" {
			DatabaseDSN = parsedEnv.OrderDatabaseDSN
		}
		if parsedEnv.LoyaltyDatabaseDSN != "" {
			LoyaltyDatabaseDSN = parsedEnv.LoyaltyDatabaseDSN
		}
		if parsedEnv.AccrualSystemAddress != "" {
			AccrualSystemAddress = parsedEnv.AccrualSystemAddress
		}
	})
}
