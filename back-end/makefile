
install:
	@echo "Installing application"
	docker-compose -f deployment/docker-compose.yaml -p lab-plataform-api up -d --build
	
cloud-install:
	@echo "Installing application in cloud"
	docker compose -f deployment/docker-compose-cloud.yaml -p lab-plataform-api up -d redis back-end --build

cloud-run:
	@echo "Running application in cloud"
	docker compose -f deployment/docker-compose-cloud.yaml -p lab-plataform-api up -d redis back-end

cloud-stop:
	@echo "Stopping application in cloud"
	docker compose -f deployment/docker-compose-cloud.yaml -p lab-plataform-api down -v --remove-orphans

stop:
	@echo "Stopping application"
	docker-compose -f deployment/docker-compose.yaml -p lab-plataform-api down -v --remove-orphans

local-install:
	@echo "Installing local setup"
	docker-compose -f ./deployment/docker-compose.yaml -p lab-plataform-api up -d redis postgres --build

local-stop:
	@echo "Stopping local setup"
	docker-compose -f ./deployment/docker-compose.yaml -p lab-plataform-api down redis postgres -v --remove-orphans

migrate-up:
	@echo "Running migrations up"
	migrate -path ./deployment/database/migrations -database ${DB_POSTGRES_URL} up

migrate-down:
	@echo "Running migrations down"
	migrate -path ./deployment/database/migrations -database ${DB_POSTGRES_URL} down

build-cloud:
	@echo "Building cloud image"
	docker buildx build --platform linux/amd64,linux/arm64 \
	-t ${DEPLOYMENT_BFF_IMAGE_NAME} \
	--push \
	-f deployment/Dockerfile \
	.
	