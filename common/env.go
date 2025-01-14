package common

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Printf("Environment variables loaded")
}

func GetString(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		Logger.Warnf("Warning: %s is not set, using %s instead", key, fallback)
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		Logger.Warnf("Warning: %s is not set, using %d instead", key, fallback)
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
		Logger.Warnf("Warning: %s is not set, using %t instead", key, fallback)
		return fallback
	}
	valAsBool, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return valAsBool
}
