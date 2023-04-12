.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: fmt vet
	go test ./... --race -coverprofile cover.out

.PHONY: repo
repo:
	cd internal && sqlc generate

.PHONY: graph
graph:
	gqlgen generate
