package common

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func EnvInit(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return fmt.Errorf("fatal: .env file not found")
	}
	return nil
}

func GetString(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Println("Warning: ", key, " is not set, using ", fallback, " instead")
		// Logger.Warnf("Warning: %s is not set, using %s instead", key, fallback)
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		log.Println("Warning: ", key, " is not set, using ", fallback, " instead")
		// Logger.Warnf("Warning: %s is not set, using %d instead", key, fallback)
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}

func GetBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		log.Println("Warning: ", key, " is not set, using ", fallback, " instead")
		// Logger.Warnf("Warning: %s is not set, using %t instead", key, fallback)
		return fallback
	}
	valAsBool, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return valAsBool
}
