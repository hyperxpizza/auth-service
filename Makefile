build:
	go build -o ./bin/main ./cmd/auth-service/main.go
docker_run:
	docker-compose up -d postgres && docker-compose up -d authservice
docker_down:
	docker-compose down --volumes
	