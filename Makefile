.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: fmt vet
	go test ./... -v --race -coverprofile cover.out

.PHONY: build
build:
	go build -o serve ./cmd/serve/main.go

.PHONY: generate
generate:
	go generate -x ./...

.PHONY: lint
lint:
	golangci-lint run -v --fix

.PHONY: multi-deploy
multi-deploy:
	weaver multi deploy weaver.toml

.PHONY: multi-status
multi-status:
	weaver multi status
