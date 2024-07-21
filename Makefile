.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: vet fmt
	~/go/bin/ginkgo -r -cover -coverprofile=coverage.out


