.PHONY: all vet fmt test build bot webapp dev-bot dev-webapp clean

all: vet test build

vet:
	@echo "Running go vet"
	@go vet ./...

test:
	@echo "Running tests"
	@go test ./...

test-cover:
	@echo "running test cover"
	@go test -v -coverprofile ./tmp/cover.out ./...
	@go tool cover -html ./tmp/cover.out -o ./tmp/cover.html
	@rm ./tmp/cover.out
	@open ./tmp/cover.html
	@sleep 1
	@rm ./tmp/cover.html

build: bot webapp

bot:
	@echo "Building bot"
	@go build -o ./bin/bot ./cmd/bot

webapp:
	@echo "Building webapp"
	@go build -o ./bin/webapp ./cmd/bot

dev-bot:
	@echo "Running bot with air for live reloading"
	@air -c .air_bot.toml 

dev-webapp:
	@echo "Running webapp with air for live reloading"
	@air -c .air_webapp.toml

clean:
	@echo "Cleaning up"
	@rm -f ./bin/bot
	@rm -f ./bin/webapp
