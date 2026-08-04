BINARY_NAME=shortlink
BUILD_DIR=build
MAIN_PATH=./cmd/main.go
ENV_PATH=configs
ENV_FILE=$(ENV_PATH)/.env

ifeq ($(OS),Windows_NT)
    RM = if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
    MKDIR = if not exist $(ENV_PATH) mkdir $(ENV_PATH)
    BINARY_WITH_EXT = $(BINARY_NAME).exe
    WRITE = echo $(1) >> $(ENV_FILE)
    CLEAR_FILE = echo. > $(ENV_FILE)
else
    RM = rm -rf $(BUILD_DIR)
    MKDIR = mkdir -p $(ENV_PATH)
    BINARY_WITH_EXT = $(BINARY_NAME)
    WRITE = echo "$(1)" >> $(ENV_FILE)
    CLEAR_FILE = printf "" > $(ENV_FILE)
endif

run:
	go run $(MAIN_PATH)

build:
	@echo "Building binary..."
	$(MKDIR)
	go build -o $(BUILD_DIR)/$(BINARY_WITH_EXT) $(MAIN_PATH)
	@echo Binary built inside $(BUILD_DIR)/$(BINARY_WITH_EXT)

.PHONY: test
test:
	@echo "Generating mocks..."
	go generate ./...
	@echo "Running tests with race detection and coverage mapping..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "Tests finished. To view coverage in browser run: go tool cover -html=coverage.out"

lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./...

clean:
	@echo "Cleaning up..."
	$(RM)
	go clean
	@echo "Cleaned!"

env:
	@echo "Generating $(ENV_FILE)..."
	@$(MKDIR)
	@echo APP_PORT=8080 > $(ENV_FILE)
	@$(call WRITE,APP_ENV=dev)
	@$(call WRITE,APP_KEY=super_secret_string_32_characters_long)
	@$(call WRITE,APP_DATABASE_URL=postgres://postgres:12345@localhost:5432/shortener?sslmode=disable)
	@$(call WRITE,APP_REDIS_URL=redis://localhost:6379/0)
	@echo ".env file created successfully in $(ENV_PATH)/"

migrate-up:
	docker build -f Dockerfile.migrate -t shortlink-migrate .
	docker run --rm --network shortlink_app-network shortlink-migrate -path=/migrations -database "postgres://postgres:12345@link_shortener_db:5432/shortener?sslmode=disable" up

migrate-down:
	docker run --rm -v //ShortLink/migrations:/migrations --network shortlink_app-network migrate/migrate -path=/migrations -database "postgres://postgres:12345@link_shortener_db:5432/shortener?sslmode=disable" down

help:
	@echo run          - run project
	@echo build        - build project
	@echo test         - start test
	@echo lint         - run golangci-lint
	@echo clean        - remove build
	@echo migrate-up   - start postgres migrations
	@echo migrate-down - stop/reverse postgres migrations
	@echo env          - make file .env
