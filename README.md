1. Проект осуществляет работу планировщика задач, позволяет пользователю добавлять, удалять и редактировать задачи. 
Задача может быть одноразовой, может иметь правило повторения.

2. Выполнены все задачи со звездочкой, кроме аутентификации

3. Локальный запуск: <go mod tidy> <go run main.go> url: <http://localhost:7540/>

4. Необходимые значения для tests/settings.go:

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = ``

5. Запуск через докер: <docker build --tag my-app:v1 .> <docker run -p 7540:7540 my-app:v1> url: <http://localhost:7540/>