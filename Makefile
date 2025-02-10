.PHONY: build clean ui

ui:
	@cd ui && pnpm pre-install && pnpm build && cd -