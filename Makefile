.PHONY: test run-server run-client

export BM_ROOT_DIR = $(shell pwd)
export BM_SERVER_CONFIG = /configs/server_config.json
export BM_CLIENT_CONFIG = /configs/client_config.json
export PORT = $(shell grep -o '"address": "[^"]*"' $(BM_ROOT_DIR)$(BM_SERVER_CONFIG) | cut -d ':' -f 3 | cut -d '"' -f 1)

export SERVER_IMAGE = book-management-server
export CLIENT_IMAGE = book-management-client

install-deps:
	go install github.com/vektra/mockery/v2

generate: install-deps
	go generate ./...

build-test-image:
	docker build -t book-management-test -f $(BM_ROOT_DIR)/deployment/test/Dockerfile .

run-tests:
	docker run --rm book-management-test

# Run tests inside a Docker container
test: build-test-image run-tests

# Build and run the server in a Docker container
run-server: stop-server
	docker network rm bm-network 2>/dev/null || true
	docker network create bm-network
	docker build -t $(SERVER_IMAGE) --build-arg BM_ROOT_DIR=$(BM_ROOT_DIR) --build-arg BM_SERVER_CONFIG=$(BM_SERVER_CONFIG) -f $(BM_ROOT_DIR)/deployment/server/Dockerfile .
	docker run -d -p $(PORT):$(PORT) --name $(SERVER_IMAGE) --network bm-network $(SERVER_IMAGE)

# Build and run the client in a Docker container
run-client: stop-client
	docker build -t $(CLIENT_IMAGE) --build-arg BM_ROOT_DIR=$(BM_ROOT_DIR) --build-arg BM_CLIENT_CONFIG=$(BM_CLIENT_CONFIG) -f $(BM_ROOT_DIR)/deployment/client/Dockerfile .
	docker run -d --name $(CLIENT_IMAGE) --network bm-network $(CLIENT_IMAGE)
	docker logs $(CLIENT_IMAGE)

# Stop and remove the server container
stop-server:
	docker stop $(SERVER_IMAGE) 2>/dev/null || true
	docker rm $(SERVER_IMAGE) 2>/dev/null || true

# Stop and remove the client container
stop-client:
	docker stop $(CLIENT_IMAGE) 2>/dev/null || true
	docker rm $(CLIENT_IMAGE) 2>/dev/null || true
