FROM golang:latest

ARG BM_ROOT_DIR
ARG BM_CLIENT_CONFIG

ENV BM_CLIENT_CONFIG=${BM_CLIENT_CONFIG}
ENV BM_ROOT_DIR=/app

WORKDIR /app

COPY . .

RUN go build -o client ./cmd/client/main.go

CMD ["./client"]
