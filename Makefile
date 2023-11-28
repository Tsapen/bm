.PHONY: test run-server run-client

export BM_ROOT_DIR = $(shell pwd)
export BM_SERVER_CONFIG = /configs/server_config.json
export BM_HTTP_CLIENT_CONFIG = /configs/http_client_config.json
export BM_CLI_CLIENT_CONFIG = /configs/cli_client_config.json
export PORT = $(shell grep -o '"address": "[^"]*"' $(BM_ROOT_DIR)$(BM_SERVER_CONFIG) | cut -d ':' -f 3 | cut -d '"' -f 1)

export SERVER_IMAGE = bm-server
export HTTP_CLIENT_IMAGE = bm-http-client
export CLI_CLIENT_IMAGE = bm-cli-client

export BM_TEST_SERVER_CONFIG = /configs/test_server_config.json
export BM_TEST_HTTP_CLIENT_CONFIG = /configs/test_server_config.json

install-deps:
	go install github.com/vektra/mockery/v2

generate: install-deps
	go generate ./...

run-tests:
	$(call print-target)
	docker-compose -f $(BM_ROOT_DIR)/deployment/test/docker-compose.yml down --remove-orphans
	docker-compose -f $(BM_ROOT_DIR)/deployment/test/docker-compose.yml up --build test-bm

test: run-tests

run-server: stop-server
	export BM_SERVER_CONFIG = /server_config.json
	docker-compose -f $(BM_ROOT_DIR)/deployment/server/docker-compose.yml down --remove-orphans
	docker-compose -f $(BM_ROOT_DIR)/deployment/server/docker-compose.yml up --build bm --detach

run-cli-client:
	docker build -t $(CLI_CLIENT_IMAGE) \
		--build-arg BM_ROOT_DIR=$(BM_ROOT_DIR) \
		--build-arg BM_CLI_CLIENT_CONFIG=$(BM_CLI_CLIENT_CONFIG) \
		--build-arg BM_CLI_CLIENT_ARGS="$(ARGS)" \
		-f $(BM_ROOT_DIR)/deployment/cli-client/Dockerfile .
	docker run $(CLI_CLIENT_IMAGE)

# Stop and remove the server container
stop-server:
	docker stop $(SERVER_IMAGE) 2>/dev/null || true
	docker rm $(SERVER_IMAGE) 2>/dev/null || true

# Stop and remove the client container
clean-cli-client:
	docker stop $(CLI_CLIENT_IMAGE) 2>/dev/null || true
	docker rm $(CLI_CLIENT_IMAGE) 2>/dev/null || true
