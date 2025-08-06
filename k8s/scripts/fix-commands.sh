#!/bin/bash

echo "=== Fix Commands Based on Analysis ==="
echo ""

echo "IMPORTANT: Based on the kustomization.yaml analysis, I see that:"
echo "1. cert-manager is NOT included in the main prod overlay"
echo "2. The prod overlay includes Google-managed certificates (managed-certificate.yaml)"
echo "3. Both cert-manager AND Google-managed certificates are configured (hybrid approach)"
echo ""

echo "=== Option 1: Apply Missing Google-Managed Certificate Resources ==="
echo "This is the quickest fix since your ingress is already configured for Google-managed certs:"
echo ""

echo "# Apply the missing resources from prod overlay"
echo "kubectl apply -f k8s/overlays/prod/backend-config.yaml"
echo "kubectl apply -f k8s/overlays/prod/managed-certificate.yaml"
echo "kubectl apply -f k8s/overlays/prod/cluster-issuer.yaml"
echo ""

echo "# OR apply the entire prod overlay (recommended)"
echo "kubectl apply -k k8s/overlays/prod/"
echo ""

echo "=== Option 2: Deploy cert-manager Separately (if needed) ==="
echo "If you want cert-manager as a backup or for other certificates:"
echo ""

echo "# Deploy cert-manager"
echo "kubectl apply -k k8s/overlays/prod/cert-manager/"
echo ""

echo "=== Option 3: Fix Static IP Assignment ==="
echo "If the ingress should use the static IP:"
echo ""

echo "# Get current static IP"
echo "STATIC_IP=\$(gcloud compute addresses describe saas-bff-jinko-test-ip --global --format='value(address)')"
echo "echo \"Static IP: \$STATIC_IP\""
echo ""

echo "# Check if ingress is using it"
echo "INGRESS_IP=\$(kubectl get ingress prod-affiliate-backend -n saas-bff -o jsonpath='{.status.loadBalancer.ingress[0].ip}')"
echo "echo \"Ingress IP: \$INGRESS_IP\""
echo ""

echo "# If they don't match, the ingress will eventually use the static IP"
echo "# You can force a refresh by deleting and recreating the ingress:"
echo "# kubectl delete ingress prod-affiliate-backend -n saas-bff"
echo "# kubectl apply -k k8s/overlays/prod/"
echo ""

echo "=== Verification Commands ==="
echo "After applying fixes, run these to verify:"
echo ""

echo "# Check certificate status"
echo "kubectl get managedcertificate affiliate-backend-ssl-cert -n saas-bff"
echo "kubectl describe managedcertificate affiliate-backend-ssl-cert -n saas-bff"
echo ""

echo "# Monitor certificate provisioning"
echo "kubectl get managedcertificate affiliate-backend-ssl-cert -n saas-bff -w"
echo ""

echo "# Test HTTPS (after certificate is Active)"
echo "curl -v https://api.affiliate.rolinko.com/health"
echo ""

echo "# Check ingress events"
echo "kubectl get events -n saas-bff --field-selector involvedObject.name=prod-affiliate-backend --sort-by='.lastTimestamp'"
echo ""

echo "=== Clean Up Migration Job ==="
echo "The migration job pod is stuck, clean it up:"
echo ""

echo "kubectl delete job prod-affiliate-backend-migration -n saas-bff"
echo ""

echo "=== Recommended Approach ==="
echo "1. Run: kubectl apply -k k8s/overlays/prod/"
echo "2. Clean up migration job: kubectl delete job prod-affiliate-backend-migration -n saas-bff"
echo "3. Monitor certificate: kubectl get managedcertificate affiliate-backend-ssl-cert -n saas-bff -w"
echo "4. Wait 10-60 minutes for certificate provisioning"
echo "5. Test HTTPS: curl https://api.affiliate.rolinko.com/health"