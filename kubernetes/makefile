install: create-cluster create-ingress-nginx setup-kubeconfig

create-cluster:
	echo "Creating cluster..."
	kind create cluster --config kind-config.yaml

create-ingress-nginx:
	echo "Creating ingress-nginx..."
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/kind/deploy.yaml
	
	echo "Labeling node kind-control-plane as ingress-ready..."
	kubectl label node kind-control-plane ingress-ready=true

setup-kubeconfig:
	echo "Exporting kubeconfig to base64..."
	cat ~/.kube/config > kubeconfig
	yq eval '.clusters = (.clusters | map(select(.name == "kind-kind").cluster.server = "https://kind-control-plane:6443"))' kubeconfig > kubeconfig_modified
	cat kubeconfig_modified | base64 -w0 > kubeconfig_base64
	rm kubeconfig
# 	grep -q '^KUBECONFIG_BASE64=' .env && \
# 		sed -i "s|^KUBECONFIG_BASE64=.*|KUBECONFIG_BASE64=$(cat kubeconfig_base64)|" .env || \
# 		echo "KUBECONFIG_BASE64=$(cat kubeconfig_base64)" >> .env
# 	rm kubeconfig_modified kubeconfig_base64

clean:
	echo "Cleaning up cluster..."
	kind delete cluster