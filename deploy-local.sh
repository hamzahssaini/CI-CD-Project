#!/bin/bash

SERVICE_NAME="go-service"
DOCKER_USER="mourphy"
IMAGE_TAG="v3"
IMAGE="$DOCKER_USER/$SERVICE_NAME:$IMAGE_TAG"
NAMESPACE="default"
HELM_RELEASE="microapp"
HELM_PATH="./helm/microapp"
VALUES_FILE="$HELM_PATH/values.yaml"

echo "🔨 Building Docker image: $IMAGE"
docker build -t $IMAGE ./microservices/$SERVICE_NAME || exit 1

echo "🚀 Pushing to Docker Hub"
docker push $IMAGE || exit 1

echo "🧠 Upgrading with Helm"
helm upgrade --install $HELM_RELEASE $HELM_PATH \
  --namespace $NAMESPACE \
  -f $VALUES_FILE \
  --set services[2].image=$IMAGE \
  --set services[2].imagePullPolicy=Always \
  --atomic --timeout 2m || exit 1

echo "♻️ Restarting pods (if needed)"
kubectl rollout restart deployment -l app=$SERVICE_NAME

echo "✅ Done!"
