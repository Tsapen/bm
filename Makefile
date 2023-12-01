export BM_ROOT_DIR = $(shell pwd)
export BM_SERVER_CONFIG = /configs/server_config.json
export BM_MIGRATIONS_PATH = /migrations/migrations.sql
export BM_HTTP_CLIENT_CONFIG = /configs/http_client_config.json
export PORT = $(shell grep -o '"address": "[^"]*"' $(BM_ROOT_DIR)$(BM_SERVER_CONFIG) | cut -d ':' -f 3 | cut -d '"' -f 1)

export SERVER_IMAGE = bm-server

export BM_TEST_SERVER_CONFIG = /configs/test_server_config.json
export BM_TEST_HTTP_CLIENT_CONFIG = /configs/test_server_config.json

test:
	$(call print-target)
	docker-compose -f $(BM_ROOT_DIR)/deployment/test/docker-compose.yml down --remove-orphans
	docker-compose -f $(BM_ROOT_DIR)/deployment/test/docker-compose.yml up -d --build test-bm

run-server:
	rm -rf $(BM_ROOT_DIR)/socket
	docker-compose -f $(BM_ROOT_DIR)/deployment/server/docker-compose.yml down --remove-orphans
	docker-compose -f $(BM_ROOT_DIR)/deployment/server/docker-compose.yml up --build bm --detach

# Stop and remove the server container
stop-server:
	docker-compose -f $(BM_ROOT_DIR)/deployment/server/docker-compose.yml down -v --remove-orphans

get-books:
	@echo "Running get-books target"; \
	SERVER_CONTAINER=$$(docker ps | grep server-bm | awk '{print $$1}'); \
	ID="$(if $(ID),--id=$(ID),)"; \
	AUTHOR="$(if $(AUTHOR),--author='$(AUTHOR)',)"; \
	GENRE="$(if $(GENRE),--genre='$(GENRE)',)"; \
	COLLECTION_ID="$(if $(COLLECTION_ID),--collection_id=$(COLLECTION_ID),)"; \
	START_DATE="$(if $(START_DATE),--start_date=$(START_DATE),)"; \
	FINISH_DATE="$(if $(FINISH_DATE),--finish_date='$(FINISH_DATE)',)"; \
	ORDER_BY="$(if $(ORDER_BY),--order_by='$(ORDER_BY)',)"; \
	DESC="$(if $(DESC),--desc=$(DESC),)"; \
	PAGE="$(if $(PAGE),--page=$(PAGE),)"; \
	PAGE_SIZE="$(if $(PAGE_SIZE),--page_size=$(PAGE_SIZE),)"; \
	docker exec -it $$SERVER_CONTAINER /bin/sh -c "./cli-client --action=get_books $$ID $$AUTHOR $$GENRE $$COLLECTION_ID $$START_DATE $$FINISH_DATE $$ORDER_BY $$DESC $$PAGE $$PAGE_SIZE"

create-book:
	@echo "Running create-book target"; \
	SERVER_CONTAINER=$$(docker ps | grep server-bm | awk '{print $$1}'); \
	TITLE="$(if $(TITLE),--title='$(TITLE)',)"; \
	AUTHOR="$(if $(AUTHOR),--author='$(AUTHOR)',)"; \
	PUBLISHED_DATE="$(if $(PUBLISHED_DATE),--published_date='$(PUBLISHED_DATE)',)"; \
	EDITION="$(if $(EDITION),--edition='$(EDITION)',)"; \
	DESCRIPTION="$(if $(DESCRIPTION),--description='$(DESCRIPTION)',)"; \
	GENRE="$(if $(GENRE),--genre='$(GENRE)',)"; \
	docker exec -it $$SERVER_CONTAINER /bin/sh -c "./cli-client --action=create_book $$TITLE $$AUTHOR $$PUBLISHED_DATE $$EDITION $$DESCRIPTION $$GENRE"

update-book:
	@echo "Running update-book target"; \
	SERVER_CONTAINER=$$(docker ps | grep server-bm | awk '{print $$1}'); \
	ID="$(if $(ID),--id=$(ID),)"; \
	TITLE="$(if $(TITLE),--title='$(TITLE)',)"; \
	AUTHOR="$(if $(AUTHOR),--author='$(AUTHOR)',)"; \
	PUBLISHED_DATE="$(if $(PUBLISHED_DATE),--published_date='$(PUBLISHED_DATE)',)"; \
	EDITION="$(if $(EDITION),--edition='$(EDITION)',)"; \
	DESCRIPTION="$(if $(DESCRIPTION),--description='$(DESCRIPTION)',)"; \
	GENRE="$(if $(GENRE),--genre='$(GENRE)',)"; \
	docker exec -it $$SERVER_CONTAINER /bin/sh -c "./cli-client --action=update_book $$ID $$TITLE $$AUTHOR $$PUBLISHED_DATE $$EDITION $$DESCRIPTION $$GENRE"

delete-books:
	@echo "Running delete-book target"; \
	SERVER_CONTAINER=$$(docker ps | grep server-bm | awk '{print $$1}'); \
	IDS="$(if $(IDS),--ids=$(IDS),)"; \
	docker exec -it $$SERVER_CONTAINER /bin/sh -c "./cli-client --action=delete_books $$IDS"

get-collections:
	@echo "Running get-collections target"; \
	SERVER_CONTAINER=$$(docker ps | grep server-bm | awk '{print $$1}'); \
	IDS="$(if $(IDS),--ids='$(IDS)',)"; \
	ORDER_BY="$(if $(ORDER_BY),--order_by='$(ORDER_BY)',)"; \
	DESC="$(if $(DESC),--desc=$(DESC),)"; \
	PAGE="$(if $(PAGE),--page=$(PAGE),)"; \
	PAGE_SIZE="$(if $(PAGE_SIZE),--page_size=$(PAGE_SIZE),)"; \
	docker exec -it $$SERVER_CONTAINER /bin/sh -c "./cli-client --action=get_collections $$IDS $$ORDER_BY $$DESC $$PAGE $$PAGE_SIZE"

create-collection:
	@echo "Running create-collection target"; \
	SERVER_CONTAINER=$$(docker ps | grep server-bm | awk '{print $$1}'); \
	NAME="$(if $(NAME),--name=$(NAME),)"; \
	DESCRIPTION="$(if $(DESCRIPTION),--description='$(DESCRIPTION)',)"; \
	docker exec -it $$SERVER_CONTAINER /bin/sh -c "./cli-client --action=create_collection $$NAME $$DESCRIPTION"

update-collection:
	@echo "Running update-collection target"; \
	SERVER_CONTAINER=$$(docker ps | grep server-bm | awk '{print $$1}'); \
	ID="$(if $(ID),--id=$(ID),)"; \
	NAME="$(if $(NAME),--name=$(NAME),)"; \
	DESCRIPTION="$(if $(DESCRIPTION),--description='$(DESCRIPTION)',)"; \
	docker exec -it $$SERVER_CONTAINER /bin/sh -c "./cli-client --action=update_collection $$ID $$NAME $$DESCRIPTION"

delete-collection:
	@echo "Running delete-collection target"; \
	SERVER_CONTAINER=$$(docker ps | grep server-bm | awk '{print $$1}'); \
	ID="$(if $(ID),--id='$(ID)',)"; \
	docker exec -it $$SERVER_CONTAINER /bin/sh -c "./cli-client --action=delete_collection $$ID"

create-books-collection:
	@echo "Running create-books-collection target"; \
	SERVER_CONTAINER=$$(docker ps | grep server-bm | awk '{print $$1}'); \
	COLLECTION_ID="$(if $(COLLECTION_ID),--collection_id='$(COLLECTION_ID)',)"; \
	BOOK_IDS="$(if $(BOOK_IDS),--book_ids='$(BOOK_IDS)',)"; \
	docker exec -it $$SERVER_CONTAINER /bin/sh -c "./cli-client --action=create_books_collection $$COLLECTION_ID $$BOOK_IDS"

delete-books-collection:
	@echo "Running delete-books-collection target"; \
	SERVER_CONTAINER=$$(docker ps | grep server-bm | awk '{print $$1}'); \
	COLLECTION_ID="$(if $(COLLECTION_ID),--collection_id='$(COLLECTION_ID)',)"; \
	BOOK_IDS="$(if $(BOOK_IDS),--book_ids='$(BOOK_IDS)',)"; \
	docker exec -it $$SERVER_CONTAINER /bin/sh -c "./cli-client --action=delete_books_collection $$COLLECTION_ID $$BOOK_IDS"
