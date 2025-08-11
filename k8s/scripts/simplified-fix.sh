#!/bin/bash

echo "=== Simplified Fix for Ingress Issues ==="
echo ""

echo "ANALYSIS:"
echo "- Static IP exists and is correctly bound to ingress (34.98.69.85)"
echo "- gcloud was querying wrong project (jinko-travel-dev vs jinko-test)"
echo "- Main issue: Missing Kubernetes resources (certificate, backend config)"
echo ""

echo "SOLUTION:"
echo "The ingress infrastructure is working correctly. We just need to apply the missing resources."
echo ""

echo "=== Commands to Run ==="
echo ""

echo "1. Fix gcloud project (if needed):"
echo "gcloud config set project jinko-test"
echo ""

echo "2. Verify static IP in correct project:"
echo "gcloud compute addresses describe saas-bff-jinko-test-ip --global --project=jinko-test"
echo ""

echo "3. Apply missing Kubernetes resources:"
echo "kubectl apply -k k8s/overlays/prod/"
echo ""

echo "4. Verify resources were created:"
echo "kubectl get managedcertificate affiliate-backend-ssl-cert -n saas-bff"
echo "kubectl get backendconfig affiliate-backend-backendconfig -n saas-bff"
echo ""

echo "5. Clean up stuck migration job:"
echo "kubectl delete job prod-affiliate-backend-migration -n saas-bff"
echo ""

echo "6. Monitor certificate provisioning (takes 10-60 minutes):"
echo "kubectl get managedcertificate affiliate-backend-ssl-cert -n saas-bff -w"
echo ""

echo "7. Test HTTPS once certificate is Active:"
echo "curl -v https://api.affiliate.rolinko.com/health"
echo ""

echo "=== Expected Timeline ==="
echo "- Immediate: Resources created, certificate starts provisioning"
echo "- 10-60 minutes: Certificate becomes Active"
echo "- HTTPS starts working"
echo ""

echo "=== Why This Will Work ==="
echo "Your infrastructure is already correctly set up:"
echo "✅ DNS points to correct IP (34.98.69.85)"
echo "✅ Static IP is bound to ingress"
echo "✅ HTTP connectivity works"
echo "✅ Application pods are healthy"
echo "✅ Service has endpoints"
echo ""
echo "❌ Missing: SSL certificate and backend configuration"
echo "   → Fixed by applying the prod overlay"