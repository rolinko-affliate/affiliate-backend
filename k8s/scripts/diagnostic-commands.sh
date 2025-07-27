#!/bin/bash

echo "=== Diagnostic Commands to Run ==="
echo "Please run these commands and share the output:"
echo ""

echo "1. Check if static IP exists and get its address:"
echo "gcloud compute addresses describe saas-bff-jinko-test-ip --global --format='value(address)'"
echo ""

echo "2. Check current ingress IP:"
echo "kubectl get ingress prod-affiliate-backend -n saas-bff -o jsonpath='{.status.loadBalancer.ingress[0].ip}'"
echo ""

echo "3. Check if ManagedCertificate exists:"
echo "kubectl get managedcertificate affiliate-backend-ssl-cert -n saas-bff"
echo ""

echo "4. Check if BackendConfig exists:"
echo "kubectl get backendconfig affiliate-backend-backendconfig -n saas-bff"
echo ""

echo "5. Check if ClusterIssuer exists:"
echo "kubectl get clusterissuer letsencrypt-prod"
echo ""

echo "6. Check what resources are actually applied in prod overlay:"
echo "kubectl kustomize k8s/overlays/prod/ | grep -E '^kind:|^  name:' | paste - -"
echo ""

echo "7. Check cert-manager resources separately:"
echo "kubectl kustomize k8s/overlays/prod/cert-manager/ | grep -E '^kind:|^  name:' | paste - -"
echo ""

echo "8. Check ingress annotations to see what it's expecting:"
echo "kubectl get ingress prod-affiliate-backend -n saas-bff -o jsonpath='{.metadata.annotations}' | jq ."
echo ""

echo "9. Check recent ingress events:"
echo "kubectl get events -n saas-bff --field-selector involvedObject.name=prod-affiliate-backend --sort-by='.lastTimestamp' | tail -10"
echo ""

echo "10. Check if there are any existing certificates or issuers:"
echo "kubectl get certificates,certificaterequests,clusterissuers,issuers --all-namespaces"
echo ""

echo "=== Analysis Questions ==="
echo "Based on the outputs above, we need to determine:"
echo "1. Are the missing resources defined but not applied?"
echo "2. Is cert-manager supposed to be deployed separately?"
echo "3. Should we use Google-managed certificates OR cert-manager (not both)?"
echo "4. Why is the static IP not being used by the ingress?"