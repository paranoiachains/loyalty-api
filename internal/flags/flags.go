package flags

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/caarlos0/env/v11"
)

// service's own flag set
var accruals = flag.NewFlagSet("", flag.ExitOnError)

// variables that were parsed from environment
var parsedEnv Environment

// parse variables only once
var once sync.Once

var (
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
)

type Environment struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func init() {
	once.Do(func() {
		accruals.StringVar(&RunAddress, "a", "localhost:8080", "service address and port")
		accruals.StringVar(&DatabaseURI, "d", "postgresql://postgres:postgres@localhost/postgres?sslmode=disable", "database connection uri")
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
		if parsedEnv.DatabaseURI != "" {
			DatabaseURI = parsedEnv.DatabaseURI
		}
		if parsedEnv.AccrualSystemAddress != "" {
			AccrualSystemAddress = parsedEnv.AccrualSystemAddress
		}
	})
}
