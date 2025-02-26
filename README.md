# CI/CD Pipeline for Go-Based Blockchain Application

This guide will walk you through setting up a **CI/CD pipeline** for your **Go-based blockchain application** that automates **testing, building, and deploying** to a local Kubernetes cluster using **Helm**.

## Prerequisites
Before proceeding, ensure you have the following installed:
- **Docker**
- **Kubernetes (k3d recommended for local clusters)**
- **Helm**
- **GitHub Actions self-hosted runner** (configured on your local machine)

## 1. Dockerize the Go Application

Create a `Dockerfile` to containerize your Go application:

```dockerfile
# Use the official Golang image as the base image
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

CMD ["./main"]
```

This multi-stage `Dockerfile` ensures a minimal final image size.

## 2. Create a Helm Chart
Helm simplifies Kubernetes deployments. Run the following command to create a Helm chart:

```bash
helm create blockchain-app
```

This command generates a directory named `blockchain-app` with a predefined structure. Customize the following files:

- **`values.yaml`**: Define the image repository and tag.
- **`templates/deployment.yaml`**: Configure Kubernetes deployments.
- **`templates/service.yaml`**: Define service exposure.

## 3. Configure GitHub Actions Workflow
Create a GitHub Actions workflow to automate testing, building, and deployment. Place the workflow file in `.github/workflows/ci-cd.yml`.

```yaml
name: CI/CD Pipeline

on:
  push:
    branches:
      - main

jobs:
  build-and-deploy:
    runs-on: self-hosted

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.20

      - name: Run unit tests
        run: go test ./...

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Log in to GitHub Container Registry
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push Docker image
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: ghcr.io/${{ github.repository_owner }}/blockchain-app:latest

      - name: Set up kubectl
        uses: azure/setup-kubectl@v3
        with:
          version: 'v1.20.0'

      - name: Set up Helm
        uses: azure/setup-helm@v3
        with:
          version: 'v3.5.0'

      - name: Deploy to Kubernetes
        run: |
          helm upgrade --install blockchain-app ./blockchain-app \
            --set image.repository=ghcr.io/${{ github.repository_owner }}/blockchain-app \
            --set image.tag=latest
```

### Key Workflow Steps:
- **Checkout Code**: Retrieves the latest repository code.
- **Set up Go**: Installs Go version 1.20.
- **Run Unit Tests**: Ensures code quality.
- **Set up Docker Buildx**: Enables multi-platform image building.
- **Authenticate to GitHub Container Registry**: Pushes the built Docker image.
- **Deploy using Helm**: Upgrades or installs the Kubernetes deployment.

## 4. Deploying to a Local k3d Cluster

Your blockchain application is deployed to a local **k3d** Kubernetes cluster. To access the application:

1. **Verify the running services**:
   ```bash
   kubectl get services
   ```
2. **Expose the service using port-forwarding**:
   ```bash
   kubectl port-forward svc/blockchain-app 8080:80
   ```
   Now, your application will be accessible at `http://localhost:8080`.

3. **Expose the service for external access using Ngrok**:
   ```bash
   ngrok http 8080
   ```
   This will provide a public URL to access your application.

## 5. Grant Kubernetes Access to GitHub Container Registry
Since Kubernetes needs access to private images, create a Kubernetes secret:

```bash
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=<your-github-username> \
  --docker-password=<your-personal-access-token> \
  --docker-email=<your-email>
```

Then, update your **`values.yaml`** to use this secret:

```yaml
imagePullSecrets:
  - name: ghcr-secret
```

## 6. Self-Hosted GitHub Runner Setup
Your **GitHub Actions self-hosted runner** is configured on your local machine to interact with the `k3d` cluster. Ensure your runner has the following dependencies installed:
- **Docker**
- **kubectl**
- **Helm**

To start the runner:
```bash
./run.sh
```

## 7. Security Considerations
- **Self-hosted runners** provide access to your local environment—restrict their usage to trusted workflows.
- **Use a GitHub PAT** with minimal required permissions.
- **Regularly update dependencies** to avoid vulnerabilities.

---
By following this guide, you will have a robust CI/CD pipeline for your **Go-based blockchain application**, ensuring automated **testing, building, and deployment** to a Kubernetes cluster. 🚀

