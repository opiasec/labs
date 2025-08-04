install: create-cluster create-registry create-ingress-nginx

create-cluster:
	echo "Creating cluster..."
	kind create cluster --config kind-config.yaml

create-registry:
	echo "Creating registry..."
	kubectl apply -f registry.yaml

create-ingress-nginx:
	echo "Creating ingress-nginx..."
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/kind/deploy.yaml
	
	echo "Labeling node kind-control-plane as ingress-ready..."
	kubectl label node kind-control-plane ingress-ready=true