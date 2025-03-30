include .env
export
.PHONY:
#Docker	
docker-build:
	docker build -t password-saver-app .
run:
	docker-compose up 
migrate:
	migrate -path pkg/db/migrations -database 'postgres://${PG_USER}:${PG_PASSWORD}@0.0.0.0:${PG_PORT}/${PG_NAME}?sslmode=disable' up

#Local Development
local-run:
	go run cmd/app/main.go
