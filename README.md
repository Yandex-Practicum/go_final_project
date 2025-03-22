# TODO-приложение

Данный проект — это веб-сервер для управления списком задач (TODO-лист), 
позволяющий добавлять, редактировать, удалять и просматривать задачи. 
Данные хранятся в SQLite.

Выполнены все задания повышенной сложности, кроме расчета даты в усложненном варианте.

### Инструкция по запуску кода локально

Для запуска потребуется установленный Go. Клонируйте репозиторий, перейдите в папку проекта и создайте файл `.env` со следующими параметрами:  
TODO_PORT=7540  
TODO_DBFILE=scheduler.db  
TODO_PASSWORD=12345

Запустите сервер командой `go run .` и откройте в браузере `http://localhost:7540`.

### Инструкция по запуску тестов

Откройте файл `tests/settings.go` и настройте параметры базы данных для тестов. 
Сейчас там все флаги установлены соответственно выполненным заданиям повышенной сложности. 
Затем выполните команду `go test ./...`.

### Инструкция по сборке и запуску через Docker

Соберите образ с помощью `docker build -t todo-app .`, затем запустите контейнер командой 
`docker run -p 7540:7540 --env-file .env --name todo-container todo-app` и откройте 
`http://localhost:7540/login.html` в браузере. В env файле внутри образа установленна только переменная окружения 
для middleware, если нужно просмотреть старые задачи, то раскомментируйте TODO_DBFILE.

Если у вас есть вопросы, пишите! 
