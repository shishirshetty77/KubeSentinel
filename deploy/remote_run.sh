#!/bin/bash
set -e

echo "🚀 Starting KubeSentinel Remote Deployment..."

# 0. Load Environment Variables (Go, etc)
export PATH=$PATH:/usr/local/go/bin:/usr/local/bin:/snap/bin

# 1. Update Repo
git pull origin master

# 2. Create Cluster (if missing)
if ! kind get clusters | grep -q "kubesentinel"; then
    echo "Creating Kind Cluster..."
    kind create cluster --name kubesentinel --config deploy/kind-config.yaml
else
    echo "Cluster 'kubesentinel' already exists."
fi

# 3. Build Docker Image
echo "Building Docker Image..."
# Use host network to speed up build if needed, or standard
docker build -t shishirshetty77/kubesentinel:latest .

# 4. Load Image into Kind
echo "Loading Image into Cluster..."
kind load docker-image shishirshetty77/kubesentinel:latest --name kubesentinel

# 5. Deploy with Helm
echo "Deploying Helm Chart..."
helm upgrade --install kubesentinel ./deploy/charts/kubesentinel \
  --set image.pullPolicy=Never \
  --set serviceAccount.create=true

# 6. Verify
echo "✅ Deployment Triggered. Waiting for pods..."
kubectl rollout status deployment/kubesentinel --timeout=60s

echo "🎉 Done! Tail the logs with:"
echo "kubectl logs -l app=kubesentinel -f"
