package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DB struct {
	Pool *pgxpool.Pool
}

func InitDB(databaseURL string) (*DB, error) {
	dbpool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	return &DB{Pool: dbpool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) ClearUsersForChat(chatID int64) error {
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM users WHERE chat_id = $1", chatID)
	return err
}

func (db *DB) AddNotification(chatID int64, message string, intervalSeconds int) error {
	log.Printf("Добавление уведомления: chatID=%d, message=%s, intervalSeconds=%d", chatID, message, intervalSeconds)
	_, err := db.Pool.Exec(context.Background(),
		"INSERT INTO notifications (chat_id, message, next_run, interval_seconds) VALUES ($1, $2, NOW() + INTERVAL '1 second' * $3, $3)",
		chatID, message, intervalSeconds)
	return err
}

func (db *DB) GetDueNotifications() ([]Notification, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT id, chat_id, message, next_run, interval_seconds FROM notifications WHERE next_run <= NOW()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.ChatID, &n.Message, &n.NextRun, &n.IntervalSeconds); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (db *DB) UpdateNotificationNextRun(id int) error {
	_, err := db.Pool.Exec(context.Background(), "UPDATE notifications SET next_run = next_run + INTERVAL '1 second' * interval_seconds WHERE id = $1", id)
	return err
}

func (db *DB) DeleteNotificationByID(notificationID int) error {
	_, err := db.Pool.Exec(context.Background(), "DELETE FROM notifications WHERE id = $1", notificationID)
	return err
}

func (db *DB) GetNotificationsForChat(chatID int64) ([]Notification, error) {
	rows, err := db.Pool.Query(context.Background(), "SELECT id, message FROM notifications WHERE chat_id = $1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.Message); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}
