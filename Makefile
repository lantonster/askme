.PHONY: pre
pre: ui orm swag wire

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