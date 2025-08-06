install: docker-compose link-network

local-install:
	docker compose up -d mongodb redis postgres

local-stop:
	docker compose down

docker-compose:
	echo "Starting Docker Compose and build images..."
	docker compose up -d --build

run:
	echo "Running Docker Compose..."
	docker compose up -d

deploy-front-end:
	echo "Deploying front-end..."
	cd front-end && npm install && npm start

deploy: install deploy-front-end

link-network:
	echo "Linking Docker Compose network..."
	docker network connect opiaseclabs_labs-net kind-control-plane

stop:
	echo "Stopping Docker Compose..."
	docker compose down

clean:
	echo "Stopping and removing Docker containers..."
	docker network disconnect opiaseclabs_labs-net kind-control-plane
	docker compose down --volumes --remove-orphans