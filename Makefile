config_path := /home/hyperxpizza/dev/golang/reusable-microservices/auth-service/config.dev.json

build:
	go build -o ./bin/main ./cmd/auth-service/main.go
docker_run:
	docker-compose up -d
docker_build_run:
	docker-compose up -d --build
docker_down:
	docker-compose down --volumes
database_test:
	go test -v ./tests --run TestInsertUser --config=$(config_path) --delete=true && \
	go test -v ./tests --run TestGetUser --config=$(config_path) --delete=true && \
	go test -v ./tests --run TestUpdateUser --config=$(config_path) --delete=true

