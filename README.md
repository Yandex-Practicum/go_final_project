# Описание проекта ---------------------
Итоговый проект по курсу "GO разработчик с нуля
Заданий со звездочкой выполнено 2, это получение данных из .env и
поиск.
Код запускается на <http://localhost:7540/>
Версия GO "go version go1.23.3 darwin/arm64"
В параметрах tests/settings.go необходимо выставить Search = true
Кроме go test -run ^TestDB$ ./tests и go test -run ^TestTasks$ ./tests
тесты не проходят не могу понять причину. Был бы очень признателен 
если научите.
Визуально приложение работает. 
ссылка:  
docker pull mozhnov/go_final_project:v1.0.3
docker run -d -p 7540:7540  mozhnov/go_final_project:v1.0.3
change 53