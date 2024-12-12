package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"time"
)

func (b *Bot) HandleMessage(msg *tgbotapi.Message) {
	if !msg.IsCommand() {
		return
	}

	switch msg.Command() {
	case "start":
		b.handleStartCommand(msg)
	case "set_notification":
		b.handleSetNotificationCommand(msg)
	case "list_notifications":
		b.handleListNotificationsCommand(msg)
	case "stop":
		b.handleStopNotificationCommand(msg)
	default:
		b.sendReply(msg.Chat.ID, "Неизвестная команда. Используйте /start, /set_notification, /list_notifications, /stop.")
	}
}

func (b *Bot) handleStartCommand(msg *tgbotapi.Message) {
	welcomeMsg := "Добро пожаловать в Notification! Используйте  /set_notification {интервал в секундах} {текст}, /list_notifications, /stop"
	b.sendReply(msg.Chat.ID, welcomeMsg)
}

func (b *Bot) handleSetNotificationCommand(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		b.sendReplyNotification(msg.Chat.ID, "Пожалуйста, укажите интервал в секундах и текст уведомления после команды /set_notification.")
		return
	}

	// Разделяем аргументы: первый - интервал, остальные - текст
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 2 {
		b.sendReplyNotification(msg.Chat.ID, "Некорректный формат. Укажите интервал и текст уведомления.")
		return
	}

	intervalSeconds, err := strconv.Atoi(parts[0])
	if err != nil || intervalSeconds <= 0 {
		b.sendReplyNotification(msg.Chat.ID, "Интервал должен быть положительным числом.")
		return
	}

	message := parts[1]

	// Логируем chat_id
	log.Printf("Добавление уведомления в группу chatID: %d, сообщение: %s", msg.Chat.ID, message)

	err = b.DB.AddNotification(msg.Chat.ID, message, intervalSeconds)
	if err != nil {
		b.sendReplyNotification(msg.Chat.ID, "Не удалось добавить уведомление. Попробуйте позже.")
		log.Printf("Ошибка при добавлении уведомления: %v", err)
		return
	}

	b.sendReplyNotification(msg.Chat.ID, fmt.Sprintf("Уведомление успешно добавлено: \"%s\" будет отправляться каждые %d секунд.", message, intervalSeconds))
}

func (b *Bot) StartNotificationScheduler() {
	go func() {
		for {
			log.Println("Проверка уведомлений...")
			notifications, err := b.DB.GetDueNotifications()
			if err != nil {
				log.Printf("Ошибка при получении уведомлений: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			log.Printf("Найдено уведомлений для отправки: %d", len(notifications)) // Логируем количество уведомлений

			for _, n := range notifications {
				log.Printf("Отправка уведомления: chatID=%d, сообщение=%s", n.ChatID, n.Message)
				b.sendReply(n.ChatID, n.Message)
				err := b.DB.UpdateNotificationNextRun(n.ID)
				if err != nil {
					log.Printf("Ошибка при обновлении времени запуска для уведомления %d: %v", n.ID, err)
				}
			}

			time.Sleep(10 * time.Second) // Проверяем уведомления каждые 10 секунд
		}
	}()
}

func (b *Bot) handleStopNotificationCommand(msg *tgbotapi.Message) {
	args := msg.CommandArguments()
	if args == "" {
		b.sendReply(msg.Chat.ID, "Пожалуйста, укажите ID уведомления для остановки. Используйте /list_notifications, чтобы увидеть доступные ID.")
		return
	}

	notificationID, err := strconv.Atoi(args)
	if err != nil {
		b.sendReply(msg.Chat.ID, "Неверный формат ID. Укажите числовой ID уведомления.")
		return
	}

	err = b.DB.DeleteNotificationByID(notificationID)
	if err != nil {
		b.sendReply(msg.Chat.ID, "Не удалось отключить уведомление. Убедитесь, что ID указан верно.")
		log.Printf("Error deleting notification with ID %d: %v", notificationID, err)
		return
	}

	b.sendReply(msg.Chat.ID, fmt.Sprintf("Уведомление с ID %d успешно отключено.", notificationID))
}

func (b *Bot) handleListNotificationsCommand(msg *tgbotapi.Message) {
	notifications, err := b.DB.GetNotificationsForChat(msg.Chat.ID)
	if err != nil {
		b.sendReply(msg.Chat.ID, "Не удалось получить список уведомлений. Попробуйте позже.")
		log.Printf("Error fetching notifications for chat %d: %v", msg.Chat.ID, err)
		return
	}

	if len(notifications) == 0 {
		b.sendReply(msg.Chat.ID, "В этом чате нет активных уведомлений.")
		return
	}

	var reply string
	for _, n := range notifications {
		reply += fmt.Sprintf("ID: %d, Текст: %s\n", n.ID, n.Message)
	}
	b.sendReply(msg.Chat.ID, "Активные уведомления:\n"+reply)
}

func (b *Bot) sendReplyNotification(chatID int64, message string) {
	log.Printf("Отправка сообщения в chatID: %d, сообщение: %s", chatID, message)
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := b.TelegramBot.Send(msg)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения в chatID %d: %v", chatID, err)
	}
}

func (b *Bot) sendReply(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := b.TelegramBot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message to chatID %d: %v", chatID, err)
	}
}
