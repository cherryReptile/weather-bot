include ${PWD}/.env

USER:=$(shell id -u)
GROUP:=$(shell id -g)

up:
	docker-compose up -d
down:
	docker-compose stop
migrate:
	docker-compose exec app migrate create -ext sql -dir migrations ${name}
migrate.up:
	docker-compose exec app migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" -path migrations up
migrate.down:
	docker-compose exec app migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):5432/$(DB_NAME)?sslmode=disable" -path migrations down