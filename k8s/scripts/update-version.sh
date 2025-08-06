#!/bin/bash
set -e

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPO_ROOT="$( cd "$SCRIPT_DIR/../.." && pwd )"

# Read version from VERSION file
VERSION=$(cat "$REPO_ROOT/VERSION")
if [ -z "$VERSION" ]; then
  echo "Error: Could not read version from VERSION file"
  exit 1
fi

echo "Updating Kubernetes manifests to use version: $VERSION"

# Update the kustomization.yaml with the current version
sed -i "s/newTag: .*/newTag: $VERSION/" "$REPO_ROOT/k8s/overlays/prod/kustomization.yaml"

# Update the migration job with the current version
sed -i "s|image: asia-east2-docker.pkg.dev/jinko-test/jinko-test-docker-repo/saas-app:.*|image: asia-east2-docker.pkg.dev/jinko-test/jinko-test-docker-repo/saas-app:$VERSION|" "$REPO_ROOT/k8s/overlays/prod/migration-job.yaml"

echo "Version updated successfully in Kubernetes manifests"