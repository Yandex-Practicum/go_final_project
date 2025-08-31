# go_final_project — Планировщик задач (TODO)

Минимальный веб-сервер на Go, реализующий TODO-лист с повторяющимися задачами.  
Хранение — SQLite. Фронтенд (HTML/CSS/JS) уже в репозитории и обслуживается сервером.

## Что делает проект

- Раздаёт фронтенд из `./web` на `http://localhost:7540/`.
- Хранит задачи в `./scheduler.db` (SQLite, таблица `scheduler`).
- REST API:
    - `POST /api/task` — добавить задачу.
    - `GET  /api/tasks` — получить ближайшие задачи (сортированы по дате).
    - `GET  /api/task?id=<id>` — получить параметры задачи.
    - `PUT  /api/task` — изменить параметры задачи.
    - `POST /api/task/done?id=<id>` — отметить выполненной (повторяемые — переносятся на следующую дату, обычные — удаляются).
    - `DELETE /api/task?id=<id>` — удалить задачу.
- Вычисляет следующую дату выполнения по правилу повторений.

## Повторения задач

Реализованы базовые и задания со звёздочкой:

- Базовые:
    - `d <N>` — через N дней (1 ≤ N ≤ 400).
    - `y` — ежегодно.
- Со звёздочкой:
    - `w <список 1..7>` — по дням недели (1 — понедельник, 7 — воскресенье).
    - `m <список дней 1..31,-1,-2> [список месяцев 1..12]` — по дням месяца с опциональным фильтром месяцев.  
      `-1` — последний день месяца, `-2` — предпоследний.

## Требования

- Go 1.25
- SQLite driver: `modernc.org/sqlite` (подтягивается `go mod tidy`).

## Запуск локально

```bash
# в корне проекта
go mod tidy
go run .
# сервер слушает http://localhost:7540/
```
## Запуск тестов
```bash
# все тесты
go test ./tests

# по шагам
go test -run ^TestApp$       ./tests   # шаг 1
go test -run ^TestDB$        ./tests   # шаг 2
go test -run ^TestNextDate$  ./tests   # шаг 3
go test -run ^TestAddTask$   ./tests   # шаг 4
go test -run ^TestTasks$     ./tests   # шаг 5
go test -run ^TestTask$      ./tests   # шаг 6 (GET)
go test -run ^TestEditTask$  ./tests   # шаг 6 (PUT)
go test -run ^TestDone$      ./tests   # шаг 7 (done)
go test -run ^TestDelTask$   ./tests   # шаг 7 (delete)

# принудительно без кеша
go test -count=1 ./tests
# или очистить кеш
go clean -testcache
```