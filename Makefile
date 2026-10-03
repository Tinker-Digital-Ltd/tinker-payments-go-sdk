.PHONY: lint test

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	elif [ -x "$${GOPATH}/bin/golangci-lint" ]; then \
		"$${GOPATH}/bin/golangci-lint" run ./...; \
	elif [ -x "$${HOME}/go/bin/golangci-lint" ]; then \
		"$${HOME}/go/bin/golangci-lint" run ./...; \
	else \
		echo "golangci-lint is required"; \
		exit 1; \
	fi

test:
	@go test ./...
