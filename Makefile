VERSION ?= latest
REGISTRY_BASE ?= 607167088920.dkr.ecr.ap-northeast-1.amazonaws.com/dreamkast-weaver

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: dbmateup
dbmateup:
	cd internal/cfp  && go run github.com/amacneil/dbmate/v2 up
	cd internal/dkui  && go run github.com/amacneil/dbmate/v2 up

.PHONY: vet
vet: dbmateup
	go vet ./...
	go run github.com/sqlc-dev/sqlc/cmd/sqlc vet -f internal/sqlc.yaml

.PHONY: test
test: fmt vet
	go test ./... -v --race -coverprofile cover.out

.PHONY: build
build:
	go build -o dkw cmd/dkw/main.go

.PHONY: generate
generate:
	go generate -x ./...
	go run github.com/ServiceWeaver/weaver/cmd/weaver generate ./...

.PHONY: lint
lint:
	golangci-lint run -v --fix

.PHONY: multi-deploy
multi-deploy:
	weaver multi deploy weaver.toml

.PHONY: multi-status
multi-status:
	weaver multi status

.PHONY: build-image
build-image:
	docker build -f Dockerfile -t $(REGISTRY_BASE):$(VERSION) .

.PHONY: push-image
push-image:
	docker push $(REGISTRY_BASE):$(VERSION)
