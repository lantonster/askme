.PHONY: pre
pre: ui orm swag wire mock

.PHONY: run
run:
	go run .

.PHONY: ui
ui:
	@cd ui && pnpm pre-install && pnpm build && cd -

.PHONY: orm
orm:
	go run cmd/orm/main.go

.PHONY: swag
swag:
	swag fmt && swag init --parseDependency

.PHONY: wire
wire:
	wire gen cmd/wire/wire.go

.PHONY: mock
mock:
	sh scripts/mock_gen.sh internal/repo mock/repo
	sh scripts/mock_gen.sh internal/service mock/service

.PHONY: test
test:
	go test $(GOFLAGS) -race -cover -coverprofile=test/cover.out ./...
	go tool cover -html=test/cover.out