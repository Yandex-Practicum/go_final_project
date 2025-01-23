
# **Todo-Go-Project**

Это веб-сервер для планирования задач с использованием SQLite. Приложение позволяет:

- Добавлять, просматривать, обновлять и удалять задачи.
- Указывать дату в формате YYYYMMDD.
- Добавлять комментарии, устанавливать правила повторения (ежегодно, каждые N дней, по конкретным дням недели/месяца) и отмечать выполнение задач.
- Фильтровать задачи по заголовку/комментарию и выбирать задачи на конкретную дату.
- Аутентифицироваться по паролю, если установлена переменная окружения `TODO_PASSWORD`, для защищенного доступа ко всем API.

**Выполненные задания:**
- Настраиваемый порт через `TODO_PORT`.
- Путь к базе данных через `TODO_DBFILE`.
- Расширенные правила повторения и возможность настраивать напоминания за 1/2 дня до конца месяца.
- Поле поиска для фильтрации задач по заголовку/комментарию или дате (формат DD.MM.YYYY).
- Аутентификация с использованием переменной окружения `TODO_PASSWORD`.


**Настройте переменные окружения (по желанию):**
- Создайте файл `.env` или используйте `env.example`:

TODO_PORT=7540
TODO_DBFILE=./database/scheduler.db
JWT_SECRET=secret_key
TODO_PASSWORD=12345


**Установите зависимости:**
go mod tidy

По умолчанию сервер будет слушать на порту 7540. Проверьте в браузере:
http://localhost:7540/

### Запуск тестов

В проекте есть пакет `tests`. Настройки хранятся в `tests/settings.go`:
package tests

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = ``


Чтобы запустить тесты, выполните:
go test ./tests

**Соберите Docker-образ:**
docker build -t todo-go .

**Запустите контейнер:**

docker run -d \
-p 7540:7540 \
-e TODO_PORT=7540 \
-e TODO_DBFILE=/app/database/scheduler.db \
-e JWT_SECRET=secret_key \
-e TODO_PASSWORD=12345 \
todo-go-app

**Проверка:**
Откройте в браузере:
http://localhost:7540/

### Примеры использования API

- **Добавление задачи:** `POST /api/tasks`
- **Просмотр задач:** `GET /api/tasks`
- **Поиск задач:** `GET /api/tasks?search=08.02.2024`
- **Обновление задачи:** `PUT /api/tasks/{id}`
- **Удаление задачи:** `DELETE /api/tasks/{id}`
- **Отметка выполнения:** `POST /api/tasks/{id}/complete`

### Заключение

`Todo-Go-Project` — веб-сервер для планирования задач на Go с использованием SQLite.
