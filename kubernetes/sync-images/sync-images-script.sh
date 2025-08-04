#!/bin/bash

set -euo pipefail

SYNC_FILE="$1"

TARGET_REGISTRY=$(grep target_registry "$SYNC_FILE" | awk '{print $2}')
IMAGES=($(grep '-' "$SYNC_FILE" | awk '{print $2}'))

echo "Registrando imagens no registry: $TARGET_REGISTRY"

# Docker Login
echo $DOCKER_PAT | docker login ghcr.io -u vitor-mauricio --password-stdin


for IMAGE in "${IMAGES[@]}"; do
    echo "🔵 Sincronizando imagem: $IMAGE"

    # Separar nome e tag
    IMAGE_NAME_WITH_TAG=$(basename "$IMAGE")
    IMAGE_NAME="${IMAGE_NAME_WITH_TAG%%:*}"
    IMAGE_TAG="${IMAGE_NAME_WITH_TAG##*:}"

    # Verificar se a imagem já existe no registry
    STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://$TARGET_REGISTRY/v2/$IMAGE_NAME/manifests/$IMAGE_TAG || true)

    if [ "$STATUS_CODE" == "200" ]; then
        echo "✅ Imagem $TARGET_REGISTRY/$IMAGE_NAME:$IMAGE_TAG já existe. Pulando push."
        continue
    fi

    echo "⚙️  Imagem não encontrada. Realizando pull, tag e push..."

    docker pull "$IMAGE"
    docker tag "$IMAGE" "$TARGET_REGISTRY/$IMAGE_NAME:$IMAGE_TAG"
    docker push "$TARGET_REGISTRY/$IMAGE_NAME:$IMAGE_TAG"
    echo "✅ Imagem $IMAGE enviada para $TARGET_REGISTRY/$IMAGE_NAME:$IMAGE_TAG"
done

echo "🎉 Sync finalizado com sucesso!"
