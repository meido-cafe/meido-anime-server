.PHONY: wire
wire:
	@cd factory && wire

.PHONY: run
run:
	@go run main.go