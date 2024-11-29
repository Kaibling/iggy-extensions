buildTime := $(shell date -u "+%Y-%m-%dT%H:%M:%S")
lint:
	golangci-lint run
fmt:
	gofumpt -l -w .
vuln:
	govulncheck ./...
gci:
	gci write --skip-generated -s standard -s default .
full-lint: gci fmt lint vuln
deps:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0

build-discord:
	cd cmd/discord; CGO_ENABLED=0 go build -ldflags "-X main.buildTime=${buildTime}" -o discord

run-discord: build-discord
	cmd/discord/discord
