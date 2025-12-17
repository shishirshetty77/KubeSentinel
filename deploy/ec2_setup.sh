#!/bin/bash
set -e

# Usage: Paste this into the "User Data" field when launching an EC2 instance
# or run it manually: `chmod +x setup.sh && ./setup.sh`

echo "🚀 Starting KubeSentinel Demonstration Environment Setup..."

# 1. Update System & Install Docker
sudo apt-get update
sudo apt-get install -y docker.io git curl wget
sudo usermod -aG docker ubuntu
newgrp docker

# 2. Install Go (Required for Kind)
export GO_VER="1.23.0"
wget https://golang.org/dl/go${GO_VER}.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go${GO_VER}.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/ubuntu/.bashrc

# 3. Install Kind (Kubernetes in Docker)
go install sigs.k8s.io/kind@v0.20.0@latest
sudo cp /home/ubuntu/go/bin/kind /usr/local/bin/

# 4. Install Kubectl
sudo snap install kubectl --classic

# 5. Install Helm
sudo snap install helm --classic

# 6. Clone Application
cd /home/ubuntu
if [ ! -d "KubeSentinel" ]; then
    git clone https://github.com/shishirshetty77/KubeSentinel.git
fi
cd KubeSentinel

echo "✅ Environment Ready! Please logout and login to refresh Docker permissions."
echo "Then run: 'make setup' (Coming soon) or follow README steps."
