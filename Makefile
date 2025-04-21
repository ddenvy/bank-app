.PHONY: init build run test docker-up docker-down migrate lint

# Полная инициализация проекта
init:
	go mod download
	go mod tidy

# Сборка приложения
build:
	go build -o bin/api cmd/api/main.go

# Запуск приложения
run:
	go run cmd/api/main.go

# Запуск тестов
test:
	go test -v ./...

# Запуск Docker контейнеров
docker-up:
	docker-compose up -d

# Остановка Docker контейнеров
docker-down:
	docker-compose down

# Применение миграций
migrate:
	docker exec -i bank_app_db psql -U postgres -d bank_app < migrations/001_init.sql

# Запуск линтера
lint:
	golangci-lint run

# Очистка
clean:
	rm -rf bin/