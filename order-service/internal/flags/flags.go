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
	AccrualSystemAddress string
)

type Environment struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseDSN          string `env:"DB_DSN"`
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
		accruals.StringVar(&DatabaseDSN, "d", "postgresql://postgres:postgres@postgres/postgres?sslmode=disable", "database connection uri")
		accruals.StringVar(&AccrualSystemAddress, "r", "", "accrual system address")
		accruals.Parse(os.Args[1:])

		err := env.Parse(&parsedEnv)
		if err != nil {
			log.Fatal(err)
		}

		// prioritize environment variables over flags
		if parsedEnv.RunAddress != "" {
			RunAddress = parsedEnv.RunAddress
		}
		if parsedEnv.DatabaseDSN != "" {
			DatabaseDSN = parsedEnv.DatabaseDSN
		}
		if parsedEnv.AccrualSystemAddress != "" {
			AccrualSystemAddress = parsedEnv.AccrualSystemAddress
		}
	})
}
