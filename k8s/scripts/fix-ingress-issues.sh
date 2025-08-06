#!/bin/bash
set -e

echo "🔧 Fixing Ingress and Certificate Issues"
echo "========================================"

NAMESPACE="saas-bff"
STATIC_IP_NAME="saas-bff-jinko-test-ip"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_step() {
    echo -e "\n${YELLOW}📋 Step $1: $2${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Step 1: Create static IP (optional - ingress already has an IP)
print_step "1" "Checking Static IP"
CURRENT_IP=$(kubectl get ingress prod-affiliate-backend -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
echo "Current ingress IP: $CURRENT_IP"

if gcloud compute addresses describe $STATIC_IP_NAME --global &>/dev/null; then
    STATIC_IP=$(gcloud compute addresses describe $STATIC_IP_NAME --global --format="value(address)")
    print_success "Static IP already exists: $STATIC_IP"
    
    if [ "$CURRENT_IP" != "$STATIC_IP" ]; then
        echo "⚠️  Warning: Ingress IP ($CURRENT_IP) differs from static IP ($STATIC_IP)"
        echo "   This is normal during initial setup. The ingress will eventually use the static IP."
    fi
else
    echo "Creating static IP..."
    if [ -n "$CURRENT_IP" ]; then
        echo "Reserving current ingress IP ($CURRENT_IP) as static IP..."
        gcloud compute addresses create $STATIC_IP_NAME --addresses=$CURRENT_IP --global
    else
        echo "Creating new static IP..."
        gcloud compute addresses create $STATIC_IP_NAME --global
    fi
    print_success "Static IP created"
fi

# Step 2: Apply missing resources
print_step "2" "Applying Missing Kubernetes Resources"

echo "Applying ClusterIssuer..."
kubectl apply -f /workspace/k8s/overlays/prod/cluster-issuer.yaml
print_success "ClusterIssuer applied"

echo "Applying BackendConfig..."
kubectl apply -f /workspace/k8s/overlays/prod/backend-config.yaml
print_success "BackendConfig applied"

echo "Applying ManagedCertificate..."
kubectl apply -f /workspace/k8s/overlays/prod/managed-certificate.yaml
print_success "ManagedCertificate applied"

# Step 3: Verify resources
print_step "3" "Verifying Applied Resources"

echo "Checking ClusterIssuer..."
if kubectl get clusterissuer letsencrypt-prod &>/dev/null; then
    print_success "ClusterIssuer exists"
else
    print_error "ClusterIssuer not found"
fi

echo "Checking BackendConfig..."
if kubectl get backendconfig affiliate-backend-backendconfig -n $NAMESPACE &>/dev/null; then
    print_success "BackendConfig exists"
else
    print_error "BackendConfig not found"
fi

echo "Checking ManagedCertificate..."
if kubectl get managedcertificate affiliate-backend-ssl-cert -n $NAMESPACE &>/dev/null; then
    print_success "ManagedCertificate exists"
    
    # Check certificate status
    CERT_STATUS=$(kubectl get managedcertificate affiliate-backend-ssl-cert -n $NAMESPACE -o jsonpath='{.status.certificateStatus}' 2>/dev/null || echo "Unknown")
    echo "Certificate status: $CERT_STATUS"
    
    if [ "$CERT_STATUS" = "Provisioning" ]; then
        echo "⏳ Certificate is being provisioned (this can take 10-60 minutes)"
    elif [ "$CERT_STATUS" = "Active" ]; then
        print_success "Certificate is active"
    else
        echo "⚠️  Certificate status: $CERT_STATUS"
    fi
else
    print_error "ManagedCertificate not found"
fi

# Step 4: Clean up migration job
print_step "4" "Cleaning Up Migration Job"
if kubectl get job prod-affiliate-backend-migration -n $NAMESPACE &>/dev/null; then
    echo "Deleting completed/failed migration job..."
    kubectl delete job prod-affiliate-backend-migration -n $NAMESPACE
    print_success "Migration job cleaned up"
else
    echo "No migration job to clean up"
fi

# Step 5: Check pod status
print_step "5" "Checking Pod Status"
echo "Current pod status:"
kubectl get pods -n $NAMESPACE -l app=affiliate-backend

READY_PODS=$(kubectl get pods -n $NAMESPACE -l app=affiliate-backend --no-headers | grep "Running" | grep "2/2" | wc -l)
TOTAL_PODS=$(kubectl get pods -n $NAMESPACE -l app=affiliate-backend --no-headers | wc -l)
echo "Ready pods: $READY_PODS/$TOTAL_PODS"

if [ $READY_PODS -eq $TOTAL_PODS ] && [ $TOTAL_PODS -gt 0 ]; then
    print_success "All application pods are ready"
else
    echo "⚠️  Some pods may not be ready. Check logs if needed:"
    echo "   kubectl logs -n $NAMESPACE -l app=affiliate-backend --tail=20"
fi

# Step 6: Test connectivity
print_step "6" "Testing Connectivity"

echo "Testing HTTP..."
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://api.affiliate.rolinko.com/health --connect-timeout 10 2>/dev/null || echo "Failed")
if [ "$HTTP_STATUS" = "200" ] || [ "$HTTP_STATUS" = "301" ] || [ "$HTTP_STATUS" = "302" ]; then
    print_success "HTTP works (status: $HTTP_STATUS)"
else
    echo "⚠️  HTTP status: $HTTP_STATUS"
fi

echo "Testing HTTPS..."
HTTPS_STATUS=$(curl -s -o /dev/null -w "%{http_code}" https://api.affiliate.rolinko.com/health --connect-timeout 10 2>/dev/null || echo "Failed")
if [ "$HTTPS_STATUS" = "200" ]; then
    print_success "HTTPS works (status: $HTTPS_STATUS)"
else
    echo "⚠️  HTTPS status: $HTTPS_STATUS (certificate may still be provisioning)"
fi

# Step 7: Monitor certificate provisioning
print_step "7" "Certificate Provisioning Status"
echo "To monitor certificate provisioning, run:"
echo "   kubectl get managedcertificate affiliate-backend-ssl-cert -n $NAMESPACE -w"
echo ""
echo "Certificate provisioning typically takes 10-60 minutes."
echo "You can also check the certificate status with:"
echo "   kubectl describe managedcertificate affiliate-backend-ssl-cert -n $NAMESPACE"

# Final summary
echo ""
echo "🎉 Setup Complete!"
echo "=================="
echo ""
echo "✅ Static IP: Created/verified"
echo "✅ ClusterIssuer: Applied"
echo "✅ BackendConfig: Applied"  
echo "✅ ManagedCertificate: Applied"
echo "✅ Migration job: Cleaned up"
echo ""
echo "Next steps:"
echo "1. Wait for certificate provisioning (10-60 minutes)"
echo "2. Monitor with: kubectl get managedcertificate affiliate-backend-ssl-cert -n $NAMESPACE -w"
echo "3. Test HTTPS once certificate is Active: curl https://api.affiliate.rolinko.com/health"
echo ""
echo "If issues persist, run the full diagnostic:"
echo "   ./k8s/scripts/check-ingress-and-certs.sh"