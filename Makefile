DAEMON_FLAG = -d
DB_USER=postgres
DB_NAME=marketplace

ifdef NO_DAEMON
    DAEMON_FLAG = 
endif

.PHONY: all
all: build up

.PHONY: build
build:
	docker-compose build

.PHONY: up
up:
	docker-compose up $(DAEMON_FLAG)

.PHONY: down
down:
	docker-compose down -v

.PHONY: register
register:
	curl -v -d '{"username":"username1", "password":"password"}' -H "Content-Type: application/json" -X POST http://localhost:8080/register

.PHONY: login
login:
	curl -v -d '{"username":"username1", "password":"password"}' -H "Content-Type: application/json" -X POST http://localhost:8080/login

.PHONY: create_ad
create_ad:
ifndef TOKEN
	$(error token undefined)
endif
	curl -v -d '{"caption":"pylesos", "description":"good pylesos", "price":123.12}' \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer $(TOKEN)" \
		-X POST http://localhost:8080/ads

.PHONY: check_ads_db
check_ads_db:
	docker exec marketplace-db psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT * FROM advertisements;"

