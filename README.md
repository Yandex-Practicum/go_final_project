# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

Сборка образа Docker:
`docker build -t final .`

Для запуска контейнера необходимо открыть терминал в текущей директории и выполнить следующие команды:
* Windows:\
`echo $null >> scheduler.db`\
`docker run -p 7540:7540 -v ${PWD}/scheduler.db:/app_go/scheduler.db final:latest`
* Unix:\
`touch scheduler.db`\
`docker run -p 7540:7540 -v $(pwd)/scheduler.db:/app_go/scheduler.db final:latest`