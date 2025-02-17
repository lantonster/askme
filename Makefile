.PHONY: ui
ui:
	@cd ui && pnpm pre-install && pnpm build && cd -

.PHONY: orm
orm:
	go run cmd/orm/main.go