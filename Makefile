APP_NAME:="bored-api-pipeline"

.PHONY: build
docker-build:
	docker build -t $(APP_NAME):latest .

.PHONY: build
build:
	CGO_ENABLED=0 go build -v -a -o $(APP_NAME) main.go

.PHONY: deploy
deploy: 
	cd helm && helm template bored-api-pipeline . --values values.yaml | kubectl apply -n bored-api -f -

.PHONY: delete
delete: 
	cd helm && helm template bored-api-pipeline . --values values.yaml | kubectl delete -n bored-api -f -

.PHONY: all
all: docker-build deploy
