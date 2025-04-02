# Файлы для итогового задания
Проект основан на учебном планировщике задач и предназначен для мониторинга, создания, редактирования и изменения статуса задач. Поддерживает правила повторения событий и хранит данные в SQLite.


# **Функциональность:**

Создание и хранение задач в базе данных SQLite.

Редактирование задач, включая изменение даты и повторяющихся правил.

Рассчет следующей даты выполнения задачи на основе правил повторения.

API для управления задачами.

Интеграция с фронтендом (директория web).

Тестирование API (директория tests).


# **Структура проекта:**

db.go – инициализация базы данных, проверка существования, создание пути и подключение к SQLite.

nextdate.go – реализация логики повторяющихся задач, проверка даты на корректность и обработка високосных годов.

tasks.go – обработчики API с проверкой входных данных.

Dockerfile – описание сборки контейнера.

go.mod / go.sum – управление зависимостями.

main.go – стартовая точка проекта, запуск веб-сервера (порт 7540), инициализация БД, маршрутизация API.

Команды для сборки
docker build -t go_final_project .
docker run -d -p 7540:7540 -v ${PWD}/scheduler.db:/app/scheduler.db go_final_project
