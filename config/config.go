package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var login = "LOGIN"
var password = "PASSWORD"
var country = "COUNTRY"
var institution = "INSTITUTION"
var apiKey = "TWO_CAPTCHA_API_KEY"

func config(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}
	return os.Getenv(key)
}

func GetLogin() string {
	return config(login)
}

func GetPassword() string {
	return config(password)
}

func GetCountry() string {
	return config(country)
}

func GetInstitution() string {
	return config(institution)
}

func GetApiKey() string {
	return config(apiKey)
}
