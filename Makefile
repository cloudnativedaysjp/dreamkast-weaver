VERSION ?= latest
REGISTRY_BASE ?= 607167088920.dkr.ecr.ap-northeast-1.amazonaws.com/dreamkast-weaver
REGISTRY_SERVE = ${REGISTRY_BASE}/serve
REGISTRY_DBMIGRATE = ${REGISTRY_BASE}/initdb

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
	weaver generate ./...

.PHONY: lint
lint:
	golangci-lint run -v --fix

.PHONY: multi-deploy
multi-deploy:
	weaver multi deploy weaver.toml

.PHONY: multi-status
multi-status:
	weaver multi status

.PHONY: initdb
initdb:
	go run tools/initdb/main.go

.PHONY: build-image
build-image:
	docker build -f dockerfiles/initdb.dockerfile -t $(REGISTRY_DBMIGRATE):$(VERSION) .
	docker build -f dockerfiles/serve.dockerfile -t $(REGISTRY_SERVE):$(VERSION) .

.PHONY: push-image
push-image:
	docker push $(REGISTRY_DBMIGRATE):$(VERSION)
	docker push $(REGISTRY_SERVE):$(VERSION)