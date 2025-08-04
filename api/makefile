
install: get-kube-config compose-up

run:
	@echo "Running docker compose"
	docker-compose -f deployments/docker-compose.yaml -p lab-cluster-api up -d

stop: compose-down

get-kube-config:
	@echo "Getting kube config"
	@cat ~/.kube/config | base64 -w 0 > kubeconfig.b64

cloud-install:
	@echo "Installing docker compose in cloud"
	docker compose -f deployments/docker-compose.yaml -p lab-cluster-api up -d --build

cloud-stop:
	@echo "Stopping docker compose in cloud"
	docker compose -f deployments/docker-compose.yaml -p lab-cluster-api down -v

cloud-run:
	@echo "Running docker compose in cloud"
	docker compose -f deployments/docker-compose.yaml -p lab-cluster-api up -d

local-install:
	@echo "Installing mongo"
	docker-compose -f deployments/docker-compose.yaml -p lab-cluster-api up -d mongodb --build

local-stop:
	@echo "Stopping mongo"
	docker-compose -f deployments/docker-compose.yaml -p lab-cluster-api down -v mongodb

compose-up:
	@echo "Starting docker compose"
	docker-compose -f deployments/docker-compose.yaml -p lab-cluster-api up -d --build

compose-down:
	@echo "Stopping docker compose"
	docker-compose -f deployments/docker-compose.yaml -p lab-cluster-api down -v

delete-all-labs-namespaces:
	@set -e
	@echo "Deleting all labs namespaces"
	@kubectl get ns -o json | jq -r '.items[].metadata.name' | grep -E '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$$' | xargs -r -t kubectl delete ns


