run:
		go run cmd/main.go --config=./internal/config/local.yaml

testall:
		go test -v ./tests/...

testallnocache:
		go test -v -count=1 ./tests/...

test1:
		go test -v -count=1 -run  TestApp ./tests

test2:
		go test -v -count=1 -run  TestDB ./tests
		
test3:
		go test -v -count=1 -run  TestNextDate ./tests
		
test4:
		go test -v -count=1 -run  TestAddTask ./tests
		
test5:
		go test -v -count=1 -run  TestTasks ./tests
		
test6:
		go test -v -count=1 -run  TestEditTask ./tests
		
test7:
		go test -v -count=1 -run  TestDone ./tests
		
test8:
		go test -v -count=1 -run  TestDelTask ./tests


docker_build:
		

docker_run:	 
		