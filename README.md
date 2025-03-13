# Go Final Project - Task Scheduler (TODO-лист)

##  Описание проекта
Этот проект — **веб-сервер на Go**, реализующий функциональность **планировщика задач (TODO-листа)**.  
Он позволяет **создавать, редактировать, удалять** и **отмечать задачи как выполненные**.  
Поддерживается **авторизация через JWT**, а также работа с базой данных **SQLite**.

**Основные функции API:**
-  **Добавление задачи**
-  **Получение списка задач** (с поиском по дате/заголовку)
-  **Редактирование задачи**
-  **Удаление задачи**
-  **Отметка задачи как выполненной**
-  **Расчет следующей даты выполнения (по правилам повторения)**

---

## Выполненные задания со звёздочкой
 **Поддержка переменных окружения (`TODO_PORT`, `TODO_DBFILE`, `TODO_PASSWORD`)**  
 **Расширенная обработка правил повторения (`d`, `y`, `w`, `m`)**  
 **Поиск по задачам (поиск по дате, заголовку, комментарию)**  
 **Аутентификация через JWT (авторизация через `/api/signin`)**  
 **Запуск через Docker (`Dockerfile` добавлен)**

---

## Инструкция по запуску проекта **локально**
### **1. Установка зависимостей**
Перед запуском убедитесь, что у вас установлен **Go (>=1.20)** и **SQLite**.  
Клонируйте репозиторий:
```sh
git clone https://github.com/yourusername/go_final_project.git
cd go_final_project
```
Загрузите зависимости:
```
go mod tidy
```
### **2. Настройка переменных окружения**
   Создайте файл .env в корневой директории и добавьте:
```
TODO_PASSWORD=yourpassword
TODO_PORT=7540
TODO_DBFILE=scheduler.db
```
(Замените yourpassword на нужный пароль.)

### **3. Запуск сервера**
```
go run cmd/server/main.go
```
После запуска откройте в браузере:
http://localhost:7540

## Инструкция по запуску тестов
Перед запуском тестов убедитесь, что сервер выключен.
Тесты используют базу ```scheduler.db```, которая может быть перезаписана.

### Настройка ```tests/settings.go```
Перед тестированием убедитесь, что переменные в ```tests/settings.go``` установлены:
```
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true  // Проверка всех правил повторения
var Search = true        // Поиск по заголовку и дате
var Token = "ваш-тестовый-токен" // Сгенерируйте токен через /api/signin
```
### Запуск всех тестов
```
go test ./tests -v
```
### Запуск конкретного теста
```
go test -run ^TestApp$ ./tests
go test -run ^TestDB$ ./tests
go test -run ^TestNextDate$ ./tests
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests
go test -run ^TestTask$ ./tests
go test -run ^TestEditTask$ ./tests
go test -run ^TestDone$ ./tests
```
## Инструкция по запуску через **Docker**
Если у вас установлен **Docker**, запустить сервер можно в контейнере.

### 1. Сборка Docker-образа
```
docker build -t my-planner .
```

### 2. Запуск контейнера
```
docker run -d -p 7540:7540 --name planner my-planner
```
Теперь откройте http://localhost:7540 в браузере.

### 3. Запуск с кастомными параметрами
```
docker run -d -p 8080:8080 \
-e TODO_PORT=8080 \
-e TODO_DBFILE=/data/scheduler.db \
--name planner my-planner
```
### 4. Остановка и удаление контейнера
```
docker stop planner && docker rm planner
```
## Структура проекта
```
go_final_project/
 ├── cmd/
 │   └── server/       
 │       └── main.go
 ├── internal/
 │   ├── auth/         # JWT-аутентификация
 │   ├── db/           # Работа с базой данных (SQLite)
 │   ├── handlers/     # API-хендлеры
 │   ├── logic/        # Бизнес-логика (повторяющиеся задачи)
 │   ├── task/         # Структура Task
 ├── web/              # Фронтенд (HTML, CSS, JS)
 ├── tests/            # Юнит-тесты
 ├── .env              # Переменные окружения
 ├── .gitignore        # Игнорируемые файлы
 ├── Dockerfile        # Запуск в Docker
 ├── go.mod / go.sum   # Зависимости Go
 ├── README.md         # Документация (этот файл)
 ├── scheduler.db      # База данных (SQLite)

```