SWAG ?= swag
GO ?= go

.PHONY: test test-cover swagger swagger-check ci-check

test:
	$(GO) test ./...

test-cover:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out

swagger:
	$(SWAG) init -g cmd/api/main.go -o docs

swagger-check:
	$(SWAG) init -g cmd/api/main.go -o docs
	git diff --exit-code -- docs/swagger.yaml docs/swagger.json docs/docs.go

ci-check: swagger-check test
