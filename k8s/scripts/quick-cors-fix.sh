#!/bin/bash

echo "🚀 Quick CORS Fix Deployment"
echo "============================"
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

echo "📦 Step 1: Building Docker image with CORS fix..."
docker build -t asia-east2-docker.pkg.dev/jinko-test/jinko-test-docker-repo/saas-app:cors-fix .

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed"
    exit 1
fi

echo "📤 Step 2: Pushing Docker image..."
docker push asia-east2-docker.pkg.dev/jinko-test/jinko-test-docker-repo/saas-app:cors-fix

if [ $? -ne 0 ]; then
    echo "❌ Docker push failed"
    exit 1
fi

echo "🔧 Step 3: Updating Kubernetes configuration..."
# Update the image tag
sed -i.bak 's/newTag: 0.0.4/newTag: cors-fix/' k8s/overlays/prod/kustomization.yaml

echo "🚀 Step 4: Deploying to Kubernetes..."
kubectl apply -k k8s/overlays/prod/

if [ $? -ne 0 ]; then
    echo "❌ Kubernetes deployment failed"
    # Restore backup
    mv k8s/overlays/prod/kustomization.yaml.bak k8s/overlays/prod/kustomization.yaml
    exit 1
fi

echo "⏳ Step 5: Waiting for rollout to complete..."
kubectl rollout status deployment/prod-affiliate-backend -n saas-bff --timeout=300s

if [ $? -ne 0 ]; then
    echo "❌ Rollout failed or timed out"
    exit 1
fi

echo ""
echo "✅ CORS fix deployed successfully!"
echo ""
echo "🧪 Test CORS with this command:"
echo "curl -X OPTIONS https://api.affiliate.rolinko.com/api/v1/users/me \\"
echo "  -H 'Origin: https://c03bbceb-ff02-4a8a-8c9c-911409c95bb8.lovableproject.com' \\"
echo "  -H 'Access-Control-Request-Method: GET' \\"
echo "  -H 'Access-Control-Request-Headers: Authorization' \\"
echo "  -v"
echo ""
echo "⚠️  REMEMBER: This is a temporary fix. Remove ALLOW_CORS_ALL after testing!"
echo ""
echo "📋 To check logs:"
echo "kubectl logs -f deployment/prod-affiliate-backend -n saas-bff -c app"