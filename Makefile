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

#Написать!!!
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

lint:
	@echo "Running golangci-lint..."
	golangci-lint run ./...

clean:
	@echo "Cleaning up..."
	@powershell -Command "if (Test-Path $(BUILD_DIR)) { Remove-Item -Recurse -Force $(BUILD_DIR) }"
	go clean
	@echo "Cleaned!"