deps:
	export GO111MODULE=on && go mod vendor

lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint sh -s -- -b $(GOBIN) latest
	GO111MODULE=on golangci-lint run ./...

test:
	GO111MODULE=on go test -mod=vendor -v -race ./...

docker:
	docker build -t corgi .
