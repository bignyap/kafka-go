package utils

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetEnvString(key string, defVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defVal
	}
	return val
}

func GetEnvInt(key string, defVal int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defVal
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Error converting %s to int: %v", key, err)
		return defVal
	}
	return intVal
}
