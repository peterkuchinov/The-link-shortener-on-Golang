BINARY_NAME=shortlink.exe
BUILD_DIR=./build
MAIN_PATH=./cmd/main.go

run:
	go run $(MAIN_PATH)

build:
	@echo "Building binary..."
	@powershell -Command "if (!(Test-Path $(BUILD_DIR))) { New-Item -ItemType Directory -Path $(BUILD_DIR) }"
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo Binary built inside $(BUILD_DIR)/$(BINARY_NAME)

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
	@powershell -Command "if (Test-Path $(BUILD_DIR)) { Remove-Item -Recurse -Force $(BUILD_DIR) }"
	go clean
	@echo "Cleaned!"

migrate-up:
	docker build -f Dockerfile.migrate -t shortlink-migrate .
	docker run --rm --network shortlink_app-network shortlink-migrate -path=/migrations -database "postgres://postgres:12345@link_shortener_db:5432/shortener?sslmode=disable" up

migrate-down:
	docker run --rm -v //ShortLink/migrations:/migrations --network shortlink_app-network migrate/migrate -path=/migrations -database "postgres://postgres:12345@link_shortener_db:5432/shortener?sslmode=disable" down



help:
	@echo run		- run project
	@echo build	 	- build project
	@echo test	 	- start test
	@echo lint	 	- run golangci-lint
	@echo clean		- remove build
	@echo migrate-up	- start postgres
	@echo migrate-down	- stop postgres

