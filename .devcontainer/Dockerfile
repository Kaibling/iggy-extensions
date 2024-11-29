FROM golang:latest
RUN gpg -k && gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69 && echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | tee /etc/apt/sources.list.d/k6.list
RUN apt update && apt install -y postgresql-client unzip k6
RUN curl -fsSL https://bun.sh/install | bash
# RUN go install github.com/go-delve/delve/cmd/dlv@latest
# RUN go install github.com/air-verse/air@latest
# RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.61.0
# RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate
# RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest