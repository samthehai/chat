#############
# VARIABLES #
#############

#############
# COMMANDS  #
#############

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

lint:
	@go mod tidy
	@golangci-lint run

lintfix:
	@golangci-lint run --fix

gqlgen:
	@echo gqlgen generating...
	rm -rf internal/interfaces/graph/generated/generated.go
	@go run github.com/99designs/gqlgen

serve:
	@go run cmd/main.go

start-redis:
	redis-server /usr/local/etc/redis.conf
