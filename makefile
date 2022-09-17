.PHONY: wire
wire:
	@cd factory && wire

.PHONY: run
run:
	@USERNAME=admin PASSWORD=admin go run main.go