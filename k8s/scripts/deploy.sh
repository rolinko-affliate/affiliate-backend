#!/bin/bash
set -e

# Default environment is prod (only prod environment available)
ENV=${1:-prod}

# Validate environment
if [[ "$ENV" != "prod" ]]; then
  echo "Error: Only 'prod' environment is available"
  echo "Usage: $0 [prod]"
  exit 1
fi

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPO_ROOT="$( cd "$SCRIPT_DIR/../.." && pwd )"

# Read version from VERSION file
VERSION=$(cat "$REPO_ROOT/VERSION")
if [ -z "$VERSION" ]; then
  echo "Error: Could not read version from VERSION file"
  exit 1
fi

echo "Deploying version $VERSION to $ENV environment..."

# Set the Docker registry and image tag
export DOCKER_REGISTRY="asia-east2-docker.pkg.dev/jinko-test/jinko-test-docker-repo/saas-app"
export IMAGE_TAG="$VERSION"

echo "Using image tag: $IMAGE_TAG"

# Update the version in Kubernetes manifests
"$SCRIPT_DIR/update-version.sh"

# Deploy the migration job first
echo "Deploying migration job..."
kubectl apply -f "$REPO_ROOT/k8s/overlays/$ENV/migration-job.yaml"

# Wait for migration to complete (optional, with timeout)
echo "Waiting for migration to complete..."
kubectl wait --for=condition=complete job/prod-affiliate-backend-migration -n affiliate-backend --timeout=300s || {
  echo "Warning: Migration job did not complete within 5 minutes. Check job status manually."
  kubectl describe job/prod-affiliate-backend-migration -n affiliate-backend
}

# Deploy the application
echo "Deploying application..."
kubectl apply -k "$REPO_ROOT/k8s/overlays/$ENV"

echo "Deployment to $ENV completed successfully!"
echo "Check deployment status with: kubectl get pods -n affiliate-backend"