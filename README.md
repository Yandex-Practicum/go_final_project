# Планировщик задач (TODO List)

Веб-приложение для управления задачами с поддержкой повторяющихся событий.

## Функциональность

- ✅ Добавление, редактирование, удаление задач
- ✅ Поддержка повторяющихся задач (ежедневно, еженедельно, ежемесячно, ежегодно)
- ✅ Отметка задач как выполненных
- ✅ Поиск задач по тексту и дате
- ✅ Аутентификация по паролю

## Выполненные задания со звёздочкой

- [x] Поддержка переменных окружения для порта и пути к БД
- [x] Поиск задач по тексту и дате
- [x] Полная реализация правил повторения (еженедельные, ежемесячные задачи)
- [x] Аутентификация с использованием JWT
- [x] Docker-образ приложения

## Локальный запуск

### Предварительные требования
- Go 1.25 или выше

### Запуск без аутентификации
```bash
go run main.go
Запуск с аутентификацией
bash
export TODO_PASSWORD="your_password"
go run main.go
Переменные окружения
TODO_PORT - порт сервера (по умолчанию: 7540)

TODO_DBFILE - путь к файлу БД (по умолчанию: scheduler.db)

TODO_PASSWORD - пароль для аутентификации (если не указан - аутентификация отключена)

Доступ к приложению
Откройте в браузере: http://localhost:7540

Запуск тестов
Без аутентификации
bash
go test ./tests
С аутентификацией
Установите пароль:

bash
export TODO_PASSWORD="testpassword"
Запустите сервер и получите токен:

bash
curl -X POST http://localhost:7540/api/signin \
  -H "Content-Type: application/json" \
  -d '{"password":"testpassword"}'
Установите токен в tests/settings.go:

go
var Token = "your_jwt_token_here"
Запустите тесты:

bash
go test ./tests
Docker
Сборка образа
bash
docker build -t todo-scheduler .
Запуск контейнера (с сохранением данных)
bash
# Создайте директорию для данных
mkdir -p docker-data

# Запустите контейнер
docker run -d \
  -p 7540:7540 \
  -v $(pwd)/docker-data:/data \
  -e TODO_PASSWORD="your_password" \
  --name todo-app \
  todo-scheduler
Управление контейнером
bash
# Просмотр логов
docker logs todo-app

# Остановка контейнера  
docker stop todo-app

# Запуск контейнера
docker start todo-app

# Удаление контейнера
docker rm todo-app
Проверка работы
После запуска откройте в браузере: http://localhost:7540

API Endpoints
Без аутентификации
GET /api/nextdate - вычисление следующей даты

POST /api/signin - аутентификация

С аутентификацией
GET /api/tasks - список задач

POST /api/task - добавление задачи

GET /api/task - получение задачи по ID

PUT /api/task - обновление задачи

DELETE /api/task - удаление задачи

POST /api/task/done - отметка задачи выполненной