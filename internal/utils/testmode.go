package utils

import (
	"os"
)

func IsTestMode() bool {
	return os.Getenv("APP_ENV") == "test"
}
