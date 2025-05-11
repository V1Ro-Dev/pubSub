COMPOSE_FILE=./deploy/docker-compose.yml
COMPOSE=docker-compose -f $(COMPOSE_FILE)

.PHONY: up
up:
	$(COMPOSE) up -d --build

.PHONY: down
down:
	$(COMPOSE) down -v

.PHONY: restart
restart: down up

.PHONY: logs
logs:
	$(COMPOSE) logs -f

.PHONY: build
build:
	$(COMPOSE) build


.PHONY: gen-proto
gen-proto:
	cd subpub/internal/delivery/grpc/proto && protoc --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto