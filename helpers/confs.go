package helpers

import (
	"fmt"
	"os"
)

func VerifyEnv(keys []string) error {
	for _, s := range keys {
		value := os.Getenv(s)
		if value == "" {
			return fmt.Errorf(s + " :: EMPTY CONFIG ENV!")
		}
	}
	return nil
}

func IsDev() bool {
	return Env() == "development"
}

func Env() string {
	return os.Getenv("ENVIRONMENT")
}
