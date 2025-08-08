## Create a Kubernetes cluster using kind and setup ingress-nginx. Create a kind cluster using kind-config.yaml for configuration
install: create-cluster create-ingress-nginx

## Create a Kubernetes cluster using kind. Using kind-config.yaml for configuration
create-cluster:
	@echo "Creating cluster..."
	@kind create cluster --config kind-config.yaml

## Setup ingress-nginx controller
create-ingress-nginx:
	@echo "Creating ingress-nginx..."
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/kind/deploy.yaml
	
	@echo "Labeling node kind-control-plane as ingress-ready..."
	@kubectl label node kind-control-plane ingress-ready=true

## Get kubeconfig and export it to base64
get-kubeconfig:
	@echo "Exporting kubeconfig to base64..."
	@cat ~/.kube/config | yq eval '.clusters = (.clusters | map(select(.name == "kind-kind").cluster.server = "https://kind-control-plane:6443"))' | base64 -w0 > kubeconfig_base64

## Clean up the cluster
clean:
	@echo "Cleaning up cluster..."
	@kind delete cluster

## Help command to show available commands
help:
	@echo "Available commands:"
	@awk '/^## / { desc = substr($$0, 4) } /^[a-zA-Z_-]+:/ && desc { printf "\033[36m%-30s\033[0m %s\n", $$1, desc; desc = "" }' $(MAKEFILE_LIST)
	