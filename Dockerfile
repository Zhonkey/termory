FROM golang:latest

RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh

WORKDIR /app

CMD ["sleep", "infinity"]