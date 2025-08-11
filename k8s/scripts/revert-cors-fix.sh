#!/bin/bash

echo "🔄 Reverting CORS Fix"
echo "===================="
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

echo "⚠️  This will remove the CORS fix and restore secure CORS settings"
echo "   Continue? (y/N)"
read -r response
if [[ ! "$response" =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 0
fi

echo "🔧 Step 1: Removing ALLOW_CORS_ALL from deployment patch..."
# Remove the CORS environment variables from deployment patch
cat > k8s/overlays/prod/patches/deployment-patch.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: affiliate-backend
  namespace: saas-bff
  labels:
    app: affiliate-backend
spec:
  replicas: 3
  template:
    spec:
      serviceAccountName: saas-bff-ksa
      containers:
      - name: app
        env:
        # Set environment to production
        - name: ENVIRONMENT
          value: "production"
        resources:
          requests:
            cpu: "200m"
            memory: "512Mi"
          limits:
            cpu: "1000m"
            memory: "1Gi"
EOF

echo "🔧 Step 2: Reverting CORS middleware..."
# Revert CORS middleware to original state
cat > internal/api/middleware/cors.go << 'EOF'
package middleware

import (
	"github.com/affiliate-backend/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a middleware that handles CORS based on the environment
func CORSMiddleware() gin.HandlerFunc {
	// In development, allow all origins
	if config.AppConfig.Environment == "development" {
		return cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		})
	}

	// In production, don't allow CORS (default behavior)
	return func(c *gin.Context) {
		c.Next()
	}
}
EOF

echo "📦 Step 3: Building secure Docker image..."
docker build -t asia-east2-docker.pkg.dev/jinko-test/jinko-test-docker-repo/saas-app:secure .

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed"
    exit 1
fi

echo "📤 Step 4: Pushing secure Docker image..."
docker push asia-east2-docker.pkg.dev/jinko-test/jinko-test-docker-repo/saas-app:secure

if [ $? -ne 0 ]; then
    echo "❌ Docker push failed"
    exit 1
fi

echo "🔧 Step 5: Updating Kubernetes configuration..."
# Update the image tag to secure version
sed -i.bak 's/newTag: cors-fix/newTag: secure/' k8s/overlays/prod/kustomization.yaml

echo "🚀 Step 6: Deploying secure version to Kubernetes..."
kubectl apply -k k8s/overlays/prod/

if [ $? -ne 0 ]; then
    echo "❌ Kubernetes deployment failed"
    # Restore backup
    mv k8s/overlays/prod/kustomization.yaml.bak k8s/overlays/prod/kustomization.yaml
    exit 1
fi

echo "⏳ Step 7: Waiting for rollout to complete..."
kubectl rollout status deployment/prod-affiliate-backend -n saas-bff --timeout=300s

if [ $? -ne 0 ]; then
    echo "❌ Rollout failed or timed out"
    exit 1
fi

echo ""
echo "✅ CORS fix reverted successfully!"
echo ""
echo "🔒 CORS is now secure (only allows same-origin requests in production)"
echo ""
echo "📋 To check logs:"
echo "kubectl logs -f deployment/prod-affiliate-backend -n saas-bff -c app"