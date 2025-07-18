package config

import (
	"fmt"
	"strconv"
	"os"

	"github.com/joho/godotenv"
)

func getEnv(key string, ifNotFound string) string {
	value, ok := os.LookupEnv(key)
	if !ok && ifNotFound == "" {
		panic(fmt.Sprintf("Missing or invalid environment key: '%s'", key))
	} else if !ok {
		return ifNotFound
	}
	return value
}

func getEnvInt(key string, defaultVal int) int {
	valStr, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		panic(fmt.Sprintf("Invalid int value for '%s': %s", key, valStr))
	}
	return val
}

func loadEnvFile(envPath string) {
	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("Error loading .env file: %s\n", envPath)
	} else {
		fmt.Printf("Successfully loaded .env file: %s\n", envPath)
	}
}
