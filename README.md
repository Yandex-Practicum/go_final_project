# Файлы для итогового задания
Проект создан на базе учебного планировщика задач, осуществляет мониторинг, создание, редактирование и изменение статуса тех или иных целей

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

Файловая структура проекта состоит из нескольких файлов, далее о них.
1. db.go - создание базы данных, проверка на существование базы, создание пути к ней, открытие подключения к SQLite
2. nextdate.go - создание правила повторения зедач, проверка на слишком старую дату и високосные года
3. tasks.go - инициалиация обработчиков с проверками
4. Dockerfile - файл для создания образа
5. go.mod и go.sum - управление зависимостями
6. main.go - стартовая точка проекта с запуском по порту 7540, инициализирует базу данных, реализует API запросы

Предполагаю что так или иначе будут ошибки, хотел бы в таком случае их обработать и закоментить непонятные моменты чтобы не перегружать код.
Буду рад обратной связи.