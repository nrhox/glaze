GOFMT?=gofmt "-s"
GOTEST=go test

.PHONY: fmt
fmt:
	@echo "Formatting & simplifying Go code..."
	@$(GOFMT) -w .
	@echo "Done!"

.PHONY: check-fmt
check-fmt:
	@echo "Checking code formatting..."
	@output=$$($(GOFMT) -l .); \
	if [ -n "$$output" ]; then \
		echo "Files not formatted:"; \
		echo "$$output"; \
		exit 1; \
	else \
		echo "All files are properly formatted"; \
	fi

.PHONY: bench-tes5
bench-test:
	$(GOTEST) -run="^\$$" -bench=. -benchmem

.PHONY: bench-test-10
bench-test-10:
	$(GOTEST) -run="^\$$" -bench=. -benchmem -benchtime=10s

.PHONY: api-test
api-test:
	$(GOTEST) -run .

.PHONY: api-test-full
api-test-full:
	$(GOTEST) ./... -v
