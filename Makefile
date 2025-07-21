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
	docker-compose down

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

.PHONY: get_ad
get_ad:
ifndef ID
	$(error ID is required. Example: make get_ad ID=123)
endif
ifndef TOKEN
	$(warning No TOKEN provided, making request without authentication)
	curl -v "http://localhost:8080/ads/$(ID)"
else
	curl -v -H "Authorization: Bearer $(TOKEN)" \
		"http://localhost:8080/ads/$(ID)"
endif

.PHONY: get_ads
get_ads:
ifndef TOKEN
	$(warning TOKEN not specified, making request without authentication)
	curl -v "http://localhost:8080/ads?page=1&page_size=10"
else
	curl -v -H "Authorization: Vearer $(TOKEN)" \
		"http://localhost:8080/ads?page=1&page_size=10"
endif

.PHONY: get_ads_filtered
get_ads_filtered:
ifndef TOKEN
	$(warning TOKEN not specified, making request without authentication)
	curl -v "http://localhost:8080/ads?page=1&page_size=10&sort_by=created_at&order=ASC"
else
	curl -v -H "Authorization: Vearer $(TOKEN)" \
		"http://localhost:8080/ads?page=1&page_size=10&sort_by=created_at&order=ASC"
endif

.PHONY: delete_ad
delete_ad:
ifndef TOKEN
	$(error TOKEN is required. Example: make delete_ad TOKEN=your_token ID=123)
endif
ifndef ID
	$(error ID is required. Example: make delete_ad TOKEN=your_token ID=123)
endif
	curl -v -H "Authorization: Bearer $(TOKEN)" \
		-X DELETE "http://localhost:8080/ads/$(ID)"

.PHONY: update_ad
update_ad:
ifndef TOKEN
	$(error TOKEN is required. Example: make update_ad TOKEN=your_token ID=123)
endif
ifndef ID
	$(error ID is required. Example: make update_ad TOKEN=your_token ID=123)
endif
	curl -v -d '{"caption":"updated", "description":"updated description", "price":200}' \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer $(TOKEN)" \
		-X PUT "http://localhost:8080/ads/$(ID)"

.PHONY: check_ads_db
check_ads_db:
	docker exec marketplace-db psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT * FROM advertisements;"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Build and start containers"
	@echo "  build        - Build containers"
	@echo "  up           - Start containers"
	@echo "  down         - Stop and remove containers"
	@echo "  register     - Register test user"
	@echo "  login        - Login test user (get TOKEN)"
	@echo "  create_ad    - Create advertisement (requires TOKEN)"
	@echo "  ads          - Get ads list (optional TOKEN)"
	@echo "  ads_filtered - Get filtered ads (optional TOKEN)"
	@echo "  update_ad    - Update advertisement (requires TOKEN and ID)"
	@echo "  delete_ad    - Delete advertisement (requires TOKEN and ID)"
	@echo "  check_ads_db - View ads in database"
	@echo ""
	@echo "Usage examples:"
	@echo "  make login"
	@echo "  make create_ad TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	@echo "  make ads TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	@echo "  make ads (without auth)"

