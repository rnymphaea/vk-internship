DAEMON_FLAG = -d

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
	curl -d '{"username":"username1", "password":"password"}' -H "Content-Type: application/json" -X POST http://localhost:8080/register
