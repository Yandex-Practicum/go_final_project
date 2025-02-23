# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

## Переменные окружения (задача со звездочками):

`TODO_PORT - порт на котором запускается веб-сервер (по умолчанию 7540)`

`TODO_DBFILE - путь к базе данных SQLite (по умолчанию ./scheduler.db)`

## Сборка Docker образа (задача со звездочкой):

`docker build -t taskmanager . && docker run taskmanager -p 7540:7540`