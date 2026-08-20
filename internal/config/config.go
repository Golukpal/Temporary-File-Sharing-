package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	MaxFileSize       int64
	AllowedFileTypes []string
}

func Load() Config {
	_ = godotenv.Load()

	maxFileSize := int64(10 * 1024 * 1024)

	if value := os.Getenv("MAX_FILE_SIZE"); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			maxFileSize = parsed
		}
	}

	allowedFileTypes := []string{
		"pdf",
		"txt",
		"png",
		"jpg",
		"jpeg",
		"zip",
	}

	if value := os.Getenv("ALLOWED_FILE_TYPES"); value != "" {
		allowedFileTypes = strings.Split(value, ",")
	}

	return Config{
		AppPort: os.Getenv("APP_PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		MaxFileSize:       maxFileSize,
		AllowedFileTypes: allowedFileTypes,
	}
}