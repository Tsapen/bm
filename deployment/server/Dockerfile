FROM golang:latest

ARG BM_ROOT_DIR
ARG BM_SERVER_CONFIG

ENV BM_ROOT_DIR=/app
ENV BM_SERVER_CONFIG=${BM_SERVER_CONFIG}
ENV BM_SERVER_CONFIG=${BM_SERVER_CONFIG}
ENV BM_MIGRATIONS_PATH=${BM_MIGRATIONS_PATH}

WORKDIR /app

COPY . .

RUN go build -o server ./cmd/server/main.go
RUN go build -o cli-client ./cmd/cli-client/main.go

CMD ["./server"]
