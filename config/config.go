// package Config. DO NOT TOUCH
package config 

import (
	"os"
	"log"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port					int `validate: "required"`
}

func LoadEnv(){
	if err :=  godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetConfig() (*Config, error){
	config := &Config{
		Port: 		getEnvAsInt("PORT", 3000),
	}

	return config, nil
}

func getEnvAsInt(name string, defaultVal int) int {
	if value, exists := os.LookupEnv(name); exists {
		return parseStringToInt(value, defaultVal)
	}
	return defaultVal
}

func parseStringToInt(s string, defaultVal int) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return value
}