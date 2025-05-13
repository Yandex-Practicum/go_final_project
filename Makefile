docker-up:
	docker compose up

docker-down:
	docker compose down -v
	
docker-start:
	docker compose start

docker-stop:
	docker compose stop

test:
	go test ./...
