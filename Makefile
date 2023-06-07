# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## Create a new confirm target
.PHONY: confirm
confirm:
	@echo -n 'Are you sure [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo 'Running an application'
	go run ./cmd/api
## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migarations/new:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migarations/up: confirm
	@echo 'Running up migarations ...'
	migrate -path ./migarateDB -database ${test.db} up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
 
## audit: tidy dependencies and format, vet and test all code 
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...