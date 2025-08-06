
## Install and setup the complete environment
install: docker-compose link-network

## Start only local databases (mongodb, redis, postgres)
local-install:
	docker compose up -d mongodb redis postgres

## Stop local Docker containers
local-stop: 
	docker compose down

## Start Docker Compose and build images
docker-compose: 
	@echo "Starting Docker Compose and build images..."
	@docker compose up -d --build

## Run Docker Compose services
run: 
	@echo "Running Docker Compose..."
	@docker compose up -d

## Deploy and start the front-end application
deploy-front-end:
	@echo "Deploying front-end..."
	@cd front-end && npm install && npm start

## Full deployment (install + front-end)
deploy: install deploy-front-end 

## Link Docker Compose network to kind cluster
link-network:
	@echo "Linking Docker Compose network..."
	@docker network connect opiaseclabs_labs-net kind-control-plane

## Stop Docker Compose Containers
stop:
	@echo "Stopping Docker Compose..."
	@docker compose stop

## Show this help message
help:
	@echo "Available commands:"
	@awk '/^## / { desc = substr($$0, 4) } /^[a-zA-Z_-]+:/ && desc { printf "\033[36m%-30s\033[0m %s\n", $$1, desc; desc = "" }' $(MAKEFILE_LIST)
	
## Stop and remove all containers, networks, and volumes
clean:
	@echo "Stopping and removing Docker containers..."
	@docker network disconnect opiaseclabs_labs-net kind-control-plane
	@docker compose down --volumes --remove-orphans