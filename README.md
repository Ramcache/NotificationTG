
# NotificationTG

## Описание
**NotificationTG** — это Telegram-бот для управления и отправки уведомлений. Он позволяет пользователям устанавливать повторяющиеся напоминания в Telegram-чатах с заданным интервалом. Проект создан на языке Go и использует PostgreSQL для хранения данных.

### Основные возможности:
- Создание уведомлений с настраиваемым интервалом.
- Список активных уведомлений.
- Удаление уведомлений.
- Планировщик для автоматической отправки сообщений.

---

## Установка и запуск

### 1. Клонирование репозитория
```bash
git clone https://github.com/Ramcache/NotificationTG
cd NotificationTG
```

### 2. Установка зависимостей
Убедитесь, что Go версии 1.16 или выше установлен. Установите зависимости:
```bash
go mod tidy
```

### 3. Настройка
Создайте файл `.env` в корневой директории и укажите параметры:
```env
TELEGRAM_TOKEN=<токен вашего Telegram-бота>
DATABASE_URL=<URL подключения к PostgreSQL>
```

### 4. Запуск
Соберите и запустите приложение:
```bash
go run main.go
```

---

## Использование

### Доступные команды:
- **`/start`** — приветственное сообщение.
- **`/set_notification <интервал> <текст>`** — создание уведомления с заданным интервалом в секундах.
- **`/list_notifications`** — просмотр списка активных уведомлений.
- **`/stop <ID>`** — удаление уведомления по ID.

---

## Структура проекта

- **`main.go`** — точка входа в приложение.
- **`internal/bot/`** — логика Telegram-бота:
  - `handler.go` — обработка команд и сообщений.
  - `bot.go` — работа с уведомлениями.
- **`internal/db/`** — работа с базой данных:
  - `db.go` — функции взаимодействия с PostgreSQL.
  - `models.go` — описание структуры данных.
- **`config/`** — загрузка и управление конфигурацией.

---

## Требования

- Go 1.16 или выше.
- PostgreSQL для хранения данных.

---

## Вклад
Мы приветствуем участие сообщества! Вы можете создать pull request или открыть issue для улучшений.

---

## Лицензия
Проект распространяется под лицензией MIT.

