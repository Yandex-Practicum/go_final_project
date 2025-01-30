# Определяем цель для сборки
build: cmd/todoapp/main.go
	go build -o todo cmd/todoapp/main.go

# Определяем цель для запуска приложения
run: todo
	./todo

# Определяем цель по умолчанию
.DEFAULT: build