# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

Проект позволяет создавать задачи и выставлять их в календарь. Есть возможность откладывать задачи на определённое время а также удалять, редактировать, отмечать выполненными. 

Все задачи со звёздочкой выполнены, но не доделана проверка с параметром var FullNextDate = true
Выдаётся ошибка, хотя алгоритм реализован. 

Почти все надписи ошибок\подсказок выполнены на английском языке, так как в моём виртуальном образе не установлен английский язык. 

Образ можно запустить вне докера с выключеным CGO=ENABLED=0 

Параметры /settings.go : 

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDI1OTcxMjQsInBhc3NIYXNoIjoiNzQxNmI2ZDQyNTI4NTk0OGIyY2RhMjk4NGU3MzgxMTkzMWQwNWYxZTExNjUxZTU2NGFkODY1NmFlZDIxOTI3YyJ9.bTXX6rj2v5d1YIDIZAnKIybgiWTTiOzXE4we2I3MoRc`
