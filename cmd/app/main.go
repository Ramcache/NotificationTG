package main

import (
	"NotificationTG/config"
	"NotificationTG/internal/bot"
	"NotificationTG/internal/db"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Проверяем, что обязательные параметры заданы
	if cfg.TelegramToken == "" || cfg.DatabaseURL == "" {
		log.Fatal("TelegramToken and DatabaseURL must be set in the configuration")
	}

	// Инициализация базы данных
	database, err := db.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer database.Close()

	// Инициализация бота
	telegramBot, err := bot.NewBot(cfg.TelegramToken, database)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}

	// Запуск планировщика уведомлений
	telegramBot.StartNotificationScheduler()

	// Запуск бота
	telegramBot.Run()
}
