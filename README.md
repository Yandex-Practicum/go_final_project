# Программа веб-сервер, который реализует функциональность простейшего планировщика задач. Это будет аналог TODO-листа

## - задания повышенной трудности не выполнялись, реализую на следующей неделе;
## - сервер запускается по адресу localhost:7540
## - В tests/settings.go следует использовать параметры по умолчанию:
var Port = 7540

var DBFile = "../scheduler.db"

var FullNextDate = false

var Search = false

var Token = ``


В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.