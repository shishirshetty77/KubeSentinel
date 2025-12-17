# KubeSentinel: The AI-Powered Kubernetes SRE

**KubeSentinel** is a custom Kubernetes Controller that autonomously detects pod failures, investigates their root cause using LLMs (Large Language Models), and suggests actionable fixes in real-time.

## 🚀 Concept

Instead of waking up at 3 AM to debug a `CrashLoopBackOff`, KubeSentinel:

1.  **Watches** for critical pod events.
2.  **Captures** logs and event context.
3.  **Analyzes** the stack trace using an AI Agent.
4.  **Reports** the root cause and fix directly to your notification channels.

## 🛠️ Architecture

- **Core**: Golang (`client-go`, `controller-runtime`)
- **Infrastructure**: Terraform (AWS EKS, VPC)
- **Packaging**: Helm Charts

## 📦 Installation

_Coming soon..._
