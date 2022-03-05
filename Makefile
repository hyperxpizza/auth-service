build:
	go build -o ./bin/main ./cmd/auth-service/main.go
docker_run:
	docker-compose up -d
docker_build_run:
	docker-compose up -d --build
docker_down:
	docker-compose down --volumes
	