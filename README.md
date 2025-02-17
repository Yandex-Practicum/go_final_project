Данный проект - это веб-сервер, написанный на языке Go, который реализует функциональность простейшего планировщика задач.

При его реализации помимо простых заданий, были выполнены дополнительные (со звездочкой):
1. Реализована возможность определять извне порт при запуске сервера.
2. Реализована возможность определять путь к файлу базы данных через переменную окружения.
3. Реализованы дополнительные правила повторения задач для планировщика.
4. Реализована возможность выбора задачи через строку поиска.
5. Реализован механизм аутентификации для доступа к планировщику по паролю.
6. Создан докер-образ с веб-сервером (написан Dockerfile). Пример командной строки для запуска контейнера таким образом, чтобы планировщик работал в браузере и подключался к SQLite базе данных на хосте, в ОС: "docker run -p 7540:7540 -e TODO_PASSWORD=12345 -v ./scheduler.db:/app/scheduler.db todo-server"

# Инструкция по запуску сервера локально
Эта инструкция поможет вам запустить сервер локально, настроить переменные окружения и подключиться к серверу через браузер.

1. Установка зависимостей
Перед запуском сервера убедитесь, что у вас установлены:

Go (версия 1.20 или выше)

SQLite3 (для работы с базой данных)

Установите зависимости проекта:

go mod download

2. Настройка переменных окружения
Создайте файл .env в корневой директории проекта. В этом файле укажите необходимые переменные окружения.

Пример .env:

TODO_PORT=7540
TODO_DBFILE=./scheduler.db
TODO_PASSWORD=12345

Описание переменных:
TODO_PORT: Порт, на котором будет запущен сервер (по умолчанию 7540).

TODO_DBFILE: Путь к файлу базы данных SQLite. Если файл не существует, он будет создан автоматически.

TODO_PASSWORD: Пароль для аутентификации на сервере.

3. Запуск сервера
Чтобы запустить сервер, выполните команду:

go run cmd/server/main.go
Сервер запустится на порту, указанном в переменной TODO_PORT (по умолчанию 7540).

4. Подключение через браузер
После запуска сервера откройте браузер и перейдите по адресу:

http://localhost:7540
Если порт был изменен в .env, укажите соответствующий порт в адресе.

5. Дополнительные флаги
Вы можете передавать дополнительные флаги при запуске сервера, чтобы переопределить переменные окружения. Например:

Указать порт:

go run cmd/server/main.go -port 8080

Указать путь к базе данных:

go run cmd/server/main.go -dbfile /path/to/scheduler.db

Указать пароль:

go run cmd/server/main.go -password mypassword

Пример команды с флагами:

go run cmd/server/main.go -port 8080 -dbfile ./my_database.db -password mypassword

6. Пример использования API
После запуска сервера вы можете взаимодействовать с ним через API. Примеры запросов:

Аутентификация:

curl -X POST http://localhost:7540/api/signin -d '{"password": "12345"}' -H "Content-Type: application/json"

Добавление задачи:

curl -X POST http://localhost:7540/api/task -d '{"title": "Новая задача", "date": "20231010"}' -H "Content-Type: application/json"

Получение списка задач:

curl http://localhost:7540/api/tasks

Обновление задачи:

curl -X PUT http://localhost:7540/api/task -d '{"id": "1", "title": "Обновленная задача"}' -H "Content-Type: application/json"

Удаление задачи:

curl -X DELETE http://localhost:7540/api/task?id=1

7. Запуск с Docker
Если вы хотите запустить сервер в Docker, используйте следующие команды:

Сборка образа:

docker build -t todo-server .

Запуск контейнера:

docker run -p 7540:7540 \
  -e TODO_PASSWORD=12345 \
  -v /путь/к/scheduler.db:/app/scheduler.db \
  todo-server
Здесь:

/путь/к/scheduler.db — путь к файлу базы данных на хосте.

TODO_PASSWORD — пароль для аутентификации.

8. Логирование и отладка
Если сервер не запускается или возникают ошибки, проверьте логи:

При локальном запуске логи выводятся в консоль.

При запуске в Docker используйте команду:

docker logs <container_id>


# Инструкция по запуску тестов: 

Перед запуском тестов укажите следующие параметры в tests/settings.go:

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = "см. примечание ниже" 

В поле Token укажите значение токена, которое сервер возвратил из /api/signin и которое хранится в куке token.

# Инструкция по сборке и запуску проекта через докер:

1. Сборка образа:

Перейдите в директорию cthdthf и выполните команду:

docker build -t todo-server .

2. Запуск контейнера:

Запустите контейнер с помощью команды:

docker run -p 7540:7540 -e TODO_PASSWORD=your_password_here todo-server

Здесь -p 7540:7540 указывает на то, что порт 7540 на хосте будет перенаправлен на порт 7540 в контейнере.

-e TODO_PASSWORD=your_password_here позволяет передать значение переменной окружения TODO_PASSWORD в контейнер.

Пример командной строки для запуска контейнера таким образом, чтобы планировщик работал в браузере и подключался к SQLite базе данных на хосте, в ОС: "docker run -p 7540:7540 -e TODO_PASSWORD=your_password_here -v ./scheduler.db:/app/scheduler.db todo-server"