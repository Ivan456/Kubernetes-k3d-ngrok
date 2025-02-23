To set up a CI/CD pipeline for your Go-based blockchain application that automates testing, building, and deploying to a local Kubernetes cluster using Helm, follow the steps below. This guide assumes you have a self-hosted GitHub Actions runner configured on your local machine, as previously discussed.

**1. Dockerize the Go Application**

Begin by creating a `Dockerfile` to containerize your Go application. This will enable consistent deployment across environments.

**`Dockerfile`**:

```dockerfile
# Use the official Golang image as the base image
FROM golang:1.20-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
```


This `Dockerfile` uses a multi-stage build to create a lightweight Docker image for your Go application.

**2. Create a Helm Chart**

Helm helps in managing Kubernetes applications. Create a Helm chart to define, install, and upgrade your Kubernetes application.

Generate a new Helm chart:


```bash
helm create blockchain-app
```


This command creates a directory named `blockchain-app` with a predefined structure. Customize the following files to suit your application:

- **`values.yaml`**: Define default configurations, such as the Docker image repository and tag.

- **`templates/deployment.yaml`**: Specify the Kubernetes deployment configuration, including container ports and environment variables.

- **`templates/service.yaml`**: Define the service that exposes your application.

**3. Configure GitHub Actions Workflow**

Set up a GitHub Actions workflow to automate testing, building, and deploying your application. Place the workflow file in `.github/workflows/ci-cd.yml` in your repository.

**`ci-cd.yml`**:

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


**Key Steps Explained**:

- **Checkout Code**: Retrieves the latest code from the repository.

- **Set up Go**: Installs the specified Go version.

- **Run Unit Tests**: Executes your Go unit tests to ensure code quality.

- **Set up Docker Buildx**: Prepares Docker Buildx for building multi-platform images.

- **Log in to GitHub Container Registry**: Authenticates to GitHub's container registry to push the Docker image.

- **Build and Push Docker Image**: Builds the Docker image and pushes it to the GitHub Container Registry.

- **Set up kubectl and Helm**: Installs `kubectl` and Helm for Kubernetes deployments.

- **Deploy to Kubernetes**: Uses Helm to deploy or upgrade the application in your local Kubernetes cluster.

**4. Ensure Kubernetes Can Pull Images from GitHub Container Registry**

Since your Kubernetes cluster needs access to the Docker images stored in the GitHub Container Registry, create a Kubernetes secret with your registry credentials:


```bash
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=<your-github-username> \
  --docker-password=<your-personal-access-token> \
  --docker-email=<your-email>
```


Replace `<your-github-username>`, `<your-personal-access-token>`, and `<your-email>` with your GitHub username, a personal access token with appropriate permissions, and your email address, respectively.

Update your Helm chart's `values.yaml` or `deployment.yaml` to use this secret for pulling images:


```yaml
imagePullSecrets:
  - name: ghcr-secret
```


**5. Security Considerations**

Running a self-hosted runner grants GitHub Actions workflows access to your local environment. Ensure that only trusted workflows are permitted to run, and regularly monitor and update your runner for security patches.

By following these steps, you can establish 