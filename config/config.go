package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DatabaseURL   string
}

func NewConfig() (*Config, error) {
	// Загружаем переменные окружения из файла .env, если он существует
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	// Считываем значения из переменных окружения
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	databaseURL := os.Getenv("DATABASE_URL")

	// Проверяем обязательные параметры
	if telegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN is required but not set")
	}
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required but not set")
	}

	return &Config{
		TelegramToken: telegramToken,
		DatabaseURL:   databaseURL,
	}, nil
}
