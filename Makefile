.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: fmt vet
	go test ./... --race -coverprofile cover.out

.PHONY: build
build:
	go build -o serve ./cmd/serve/main.go

.PHONY: generate
generate:
	go generate ./...

.PHONY: multi-deploy
multi-deploy:
	weaver multi deploy weaver.toml

.PHONY: multi-status
multi-status:
	weaver multi status
