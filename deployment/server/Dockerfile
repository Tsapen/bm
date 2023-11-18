FROM golang:latest

ARG BM_ROOT_DIR
ARG BM_SERVER_CONFIG

ENV BM_SERVER_CONFIG=${BM_SERVER_CONFIG}
ENV BM_ROOT_DIR=/app

WORKDIR /app

COPY . .

RUN go build -o server ./cmd/server/main.go

CMD ["./server"]