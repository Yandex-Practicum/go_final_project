# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

-Данный код представляет из себя веб-сервер, который реализует функциональность простейшего планировщика задач. 
API содержит следующие операции:
добавить задачу;
получить список задач;
удалить задачу;
получить параметры задачи;
изменить параметры задачи;
отметить задачу как выполненную.

В файле `main.go `
Реализована функция `createTable`, создающая таблицу и `taskHandler` обработчик для GET, POST, Put запросов к /api/task
В файле `RulesNextDate.go` реализована функция `NextDate` представляющую из себя правила повторения задач. (Каждые d дней или же ежегодно) 
В файле `task.go` реализованы обработчики API запросов на добавление, получение , удаление и другие., задач.



-Задания со звездочкой не выполнялись

-Проект запускается через project.exe, или через консоль командой "go run .". на сайт в браузере можно перейти по http://localhost:7540/

-Запустить тесты можер одной командой go test ./tests

