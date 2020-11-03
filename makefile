all:
	gofmt -w main.go
	docker-compose stop
	docker-compose build
	docker-compose up
