include ${PWD}/.env

USER:=$(shell id -u)
GROUP:=$(shell id -g)

build:
	docker compose run --rm app sh -c "CGO_ENABLED=0 go build -o tmp/app cmd/main.go"
run:
	docker compose run --rm app sh -c "CGO_ENABLED=0 go build -o tmp/app cmd/main.go && ./tmp/app"
up:
	docker-compose up -d && make log
down:
	docker-compose stop
exec:
	docker-compose exec app sh
migrate:
	docker-compose exec app migrate create -ext sql -dir migrations ${name}
migrate.up:
	docker-compose exec app migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" -path migrations up
migrate.down:
	docker-compose exec app migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" -path migrations down
log:
	docker-compose logs -f app