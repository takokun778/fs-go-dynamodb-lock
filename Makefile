export 
CONTAINER_NAME := dynamo
DYNAMODB_PORT := 8888
DYNAMODB_URL := http://localhost:$(DYNAMODB_PORT)


.PHONY: db

db:
	@docker run --rm -d \
		-p $(DYNAMODB_PORT):8000 \
		--name $(CONTAINER_NAME) \
		amazon/dynamodb-local:1.18.0

stop:
	@docker stop $(CONTAINER_NAME)

run:
	@go run main.go
