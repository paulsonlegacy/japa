package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
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

func loadEnvFile(file string) {
	if err := godotenv.Load(file); err != nil {
		fmt.Printf("Error on load environment file: %s", file)
	}
}
