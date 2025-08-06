#!/bin/bash

echo "=== Corrected Diagnostic Commands ==="
echo ""

echo "ISSUE IDENTIFIED: gcloud is using wrong project!"
echo "Your cluster is in 'jinko-test' but gcloud is querying 'jinko-travel-dev'"
echo ""

echo "1. Check current gcloud project:"
echo "gcloud config get-value project"
echo ""

echo "2. Set correct project (if needed):"
echo "gcloud config set project jinko-test"
echo ""

echo "3. Now check static IP in correct project:"
echo "gcloud compute addresses describe saas-bff-jinko-test-ip --global --format='value(address)' --project=jinko-test"
echo ""

echo "4. Alternative - list all static IPs to find it:"
echo "gcloud compute addresses list --global --project=jinko-test"
echo ""

echo "5. Check if the current ingress IP (34.98.69.85) is the static IP:"
echo "gcloud compute addresses list --global --filter='address:34.98.69.85' --project=jinko-test"
echo ""

echo "6. Check what resources would be created by prod overlay:"
echo "kubectl kustomize k8s/overlays/prod/ | grep -E '^kind:|^  name:' | paste - -"
echo ""

echo "7. Check if missing resources exist:"
echo "kubectl get managedcertificate affiliate-backend-ssl-cert -n saas-bff"
echo "kubectl get backendconfig affiliate-backend-backendconfig -n saas-bff"
echo "kubectl get clusterissuer letsencrypt-prod"
echo ""

echo "8. Check ingress annotations:"
echo "kubectl get ingress prod-affiliate-backend -n saas-bff -o jsonpath='{.metadata.annotations}' | jq ."
echo ""

echo "=== Key Insight ==="
echo "If the ingress IP (34.98.69.85) matches your static IP, then everything is working correctly!"
echo "The main issue is likely just the missing certificate and backend config resources."
echo ""

echo "=== Quick Fix Commands ==="
echo "# Apply missing resources"
echo "kubectl apply -k k8s/overlays/prod/"
echo ""
echo "# Clean up migration job"
echo "kubectl delete job prod-affiliate-backend-migration -n saas-bff"
echo ""
echo "# Monitor certificate"
echo "kubectl get managedcertificate affiliate-backend-ssl-cert -n saas-bff -w"