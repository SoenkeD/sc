.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: vet fmt
	~/go/bin/ginkgo -r -cover -coverprofile=coverage.out


COMMIT_HASH=$(shell git rev-parse HEAD)
.PHONY: build
build:
	go build -ldflags "-X main.commitHash=$(COMMIT_HASH)"


