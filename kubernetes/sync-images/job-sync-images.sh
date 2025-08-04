#!/bin/bash

set -euo pipefail

# Parâmetros
DOCKER_PAT=$1

# Variáveis
NAMESPACE="registry"
CONFIGMAP_NAME="sync-images-file"
CONFIGMAP_SCRIPT_NAME="sync-images-script"
JOB_NAME="sync-images-job"
SYNC_FILE_PATH="./sync-images.yaml"
SYNC_SCRIPT_FILE_PATH="./sync-images-script.sh"

# 1. Criar ConfigMap temporário
echo "🔵 Criando ConfigMap com o sync-images.yaml..."
kubectl create configmap $CONFIGMAP_NAME --from-file=sync-images.yaml=$SYNC_FILE_PATH -n $NAMESPACE

echo "🔵 Criando ConfigMap com o sync-images-script.sh..."
kubectl create configmap $CONFIGMAP_SCRIPT_NAME --from-file=sync-images-script.sh=$SYNC_SCRIPT_FILE_PATH -n $NAMESPACE

# 2. Criar Secret
echo "🔵 Criando Secret com o Docker PAT..."
kubectl create secret generic registry-credentials --from-literal=DOCKER_PAT=$DOCKER_PAT -n $NAMESPACE

# 3. Aplicar o Job
echo "🚀 Aplicando o Job de sincronização..."
kubectl apply -f job-sync-images.yaml

# 4. Aguardar o Job completar
echo "⏳ Aguardando o Job terminar..."
kubectl wait --for=condition=complete job/$JOB_NAME -n $NAMESPACE --timeout=300s

# 5. Limpar Configs
echo "🧹 Limpando Configs..."
kubectl delete configmap $CONFIGMAP_NAME -n $NAMESPACE
kubectl delete configmap $CONFIGMAP_SCRIPT_NAME -n $NAMESPACE
kubectl delete secret registry-credentials -n $NAMESPACE

echo "✅ Sincronização concluída e limpeza feita!"
