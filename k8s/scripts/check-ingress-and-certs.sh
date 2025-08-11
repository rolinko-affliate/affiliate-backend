#!/bin/bash
set -e

echo "=== Ingress and Certificate Manager Setup Checker ==="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables
NAMESPACE="saas-bff"
INGRESS_NAME="prod-affiliate-backend"
CERT_NAME="affiliate-backend-ssl-cert"
DOMAIN="api.affiliate.rolinko.com"
STATIC_IP_NAME="saas-bff-jinko-test-ip"

# Helper functions
print_header() {
    echo -e "\n${BLUE}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

check_command() {
    if command -v $1 &> /dev/null; then
        print_success "$1 is installed"
        return 0
    else
        print_error "$1 is not installed"
        return 1
    fi
}

# Check prerequisites
print_header "Prerequisites Check"
check_command kubectl || exit 1
check_command gcloud || print_warning "gcloud CLI not found - some GCP-specific checks will be skipped"

# Check if we can connect to cluster
if kubectl cluster-info &> /dev/null; then
    print_success "Connected to Kubernetes cluster"
    CLUSTER_NAME=$(kubectl config current-context)
    echo "Current context: $CLUSTER_NAME"
else
    print_error "Cannot connect to Kubernetes cluster"
    exit 1
fi

# Check namespace
print_header "Namespace Check"
if kubectl get namespace $NAMESPACE &> /dev/null; then
    print_success "Namespace '$NAMESPACE' exists"
else
    print_error "Namespace '$NAMESPACE' does not exist"
    echo "Create it with: kubectl create namespace $NAMESPACE"
fi

# Check cert-manager installation
print_header "Cert-Manager Installation Check"
if kubectl get namespace cert-manager &> /dev/null; then
    print_success "cert-manager namespace exists"
    
    # Check cert-manager pods
    echo "Cert-manager pods status:"
    kubectl get pods -n cert-manager -o wide
    
    # Check if all cert-manager pods are ready
    NOT_READY=$(kubectl get pods -n cert-manager --no-headers | grep -v "Running\|Completed" | wc -l)
    if [ $NOT_READY -eq 0 ]; then
        print_success "All cert-manager pods are running"
    else
        print_warning "$NOT_READY cert-manager pods are not ready"
    fi
    
    # Check cert-manager version
    CERT_MANAGER_VERSION=$(kubectl get deployment cert-manager -n cert-manager -o jsonpath='{.spec.template.spec.containers[0].image}' | cut -d':' -f2)
    echo "Cert-manager version: $CERT_MANAGER_VERSION"
    
else
    print_error "cert-manager is not installed"
    echo "Install it by applying: kubectl apply -f k8s/base/cert-manager/"
fi

# Check ClusterIssuer
print_header "ClusterIssuer Check"
if kubectl get clusterissuer letsencrypt-prod &> /dev/null; then
    print_success "ClusterIssuer 'letsencrypt-prod' exists"
    
    # Check ClusterIssuer status
    ISSUER_STATUS=$(kubectl get clusterissuer letsencrypt-prod -o jsonpath='{.status.conditions[0].status}' 2>/dev/null || echo "Unknown")
    if [ "$ISSUER_STATUS" = "True" ]; then
        print_success "ClusterIssuer is ready"
    else
        print_warning "ClusterIssuer status: $ISSUER_STATUS"
        echo "ClusterIssuer details:"
        kubectl describe clusterissuer letsencrypt-prod
    fi
else
    print_error "ClusterIssuer 'letsencrypt-prod' does not exist"
    echo "Apply it with: kubectl apply -f k8s/overlays/prod/cluster-issuer.yaml"
fi

# Check static IP (GCP specific)
print_header "Static IP Check (GCP)"
if command -v gcloud &> /dev/null; then
    if gcloud compute addresses describe $STATIC_IP_NAME --global &> /dev/null; then
        STATIC_IP=$(gcloud compute addresses describe $STATIC_IP_NAME --global --format="value(address)")
        print_success "Static IP '$STATIC_IP_NAME' exists: $STATIC_IP"
    else
        print_error "Static IP '$STATIC_IP_NAME' does not exist"
        echo "Create it with: gcloud compute addresses create $STATIC_IP_NAME --global"
    fi
else
    print_warning "gcloud CLI not available - skipping static IP check"
fi

# Check ingress
print_header "Ingress Check"
if kubectl get ingress $INGRESS_NAME -n $NAMESPACE &> /dev/null; then
    print_success "Ingress '$INGRESS_NAME' exists"
    
    echo "Ingress details:"
    kubectl get ingress $INGRESS_NAME -n $NAMESPACE -o wide
    
    # Check ingress IP
    INGRESS_IP=$(kubectl get ingress $INGRESS_NAME -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
    if [ -n "$INGRESS_IP" ]; then
        print_success "Ingress has IP address: $INGRESS_IP"
    else
        print_warning "Ingress does not have an IP address yet"
    fi
    
    # Check ingress annotations
    echo "Ingress annotations:"
    kubectl get ingress $INGRESS_NAME -n $NAMESPACE -o jsonpath='{.metadata.annotations}' | jq . 2>/dev/null || kubectl get ingress $INGRESS_NAME -n $NAMESPACE -o jsonpath='{.metadata.annotations}'
    
else
    print_error "Ingress '$INGRESS_NAME' does not exist"
    echo "Apply it with: kubectl apply -k k8s/overlays/prod/"
fi

# Check managed certificate (GKE specific)
print_header "Managed Certificate Check (GKE)"
if kubectl get managedcertificate $CERT_NAME -n $NAMESPACE &> /dev/null; then
    print_success "ManagedCertificate '$CERT_NAME' exists"
    
    # Check certificate status
    CERT_STATUS=$(kubectl get managedcertificate $CERT_NAME -n $NAMESPACE -o jsonpath='{.status.certificateStatus}' 2>/dev/null || echo "Unknown")
    echo "Certificate status: $CERT_STATUS"
    
    if [ "$CERT_STATUS" = "Active" ]; then
        print_success "Certificate is active"
    else
        print_warning "Certificate is not active yet (status: $CERT_STATUS)"
        echo "Certificate details:"
        kubectl describe managedcertificate $CERT_NAME -n $NAMESPACE
    fi
else
    print_error "ManagedCertificate '$CERT_NAME' does not exist"
    echo "Apply it with: kubectl apply -f k8s/overlays/prod/managed-certificate.yaml"
fi

# Check backend config
print_header "BackendConfig Check"
BACKEND_CONFIG_NAME="affiliate-backend-backendconfig"
if kubectl get backendconfig $BACKEND_CONFIG_NAME -n $NAMESPACE &> /dev/null; then
    print_success "BackendConfig '$BACKEND_CONFIG_NAME' exists"
    
    echo "BackendConfig details:"
    kubectl get backendconfig $BACKEND_CONFIG_NAME -n $NAMESPACE -o yaml
else
    print_error "BackendConfig '$BACKEND_CONFIG_NAME' does not exist"
    echo "Apply it with: kubectl apply -f k8s/overlays/prod/backend-config.yaml"
fi

# Check service
print_header "Service Check"
SERVICE_NAME="prod-affiliate-backend"
if kubectl get service $SERVICE_NAME -n $NAMESPACE &> /dev/null; then
    print_success "Service '$SERVICE_NAME' exists"
    
    echo "Service details:"
    kubectl get service $SERVICE_NAME -n $NAMESPACE -o wide
    
    # Check service endpoints
    ENDPOINTS=$(kubectl get endpoints $SERVICE_NAME -n $NAMESPACE -o jsonpath='{.subsets[*].addresses[*].ip}' 2>/dev/null || echo "")
    if [ -n "$ENDPOINTS" ]; then
        print_success "Service has endpoints: $ENDPOINTS"
    else
        print_warning "Service has no endpoints - check if pods are running"
    fi
else
    print_error "Service '$SERVICE_NAME' does not exist"
fi

# Check pods
print_header "Pod Check"
PODS=$(kubectl get pods -n $NAMESPACE -l app=affiliate-backend --no-headers 2>/dev/null | wc -l)
if [ $PODS -gt 0 ]; then
    print_success "$PODS pod(s) found"
    kubectl get pods -n $NAMESPACE -l app=affiliate-backend -o wide
    
    # Check pod readiness
    READY_PODS=$(kubectl get pods -n $NAMESPACE -l app=affiliate-backend --no-headers | grep "Running" | grep "1/1" | wc -l)
    if [ $READY_PODS -eq $PODS ]; then
        print_success "All pods are ready"
    else
        print_warning "Only $READY_PODS out of $PODS pods are ready"
    fi
else
    print_error "No pods found with label app=affiliate-backend"
fi

# DNS and connectivity check
print_header "DNS and Connectivity Check"
echo "Checking DNS resolution for $DOMAIN..."
if nslookup $DOMAIN &> /dev/null; then
    RESOLVED_IP=$(nslookup $DOMAIN | grep "Address:" | tail -1 | awk '{print $2}')
    print_success "DNS resolves $DOMAIN to $RESOLVED_IP"
    
    # Compare with ingress IP
    if [ -n "$INGRESS_IP" ] && [ "$RESOLVED_IP" = "$INGRESS_IP" ]; then
        print_success "DNS IP matches ingress IP"
    elif [ -n "$INGRESS_IP" ]; then
        print_warning "DNS IP ($RESOLVED_IP) does not match ingress IP ($INGRESS_IP)"
    fi
else
    print_warning "DNS resolution failed for $DOMAIN"
fi

# HTTP/HTTPS connectivity check
echo "Testing HTTP connectivity..."
if curl -s -o /dev/null -w "%{http_code}" http://$DOMAIN/health --connect-timeout 10 | grep -q "200\|301\|302"; then
    print_success "HTTP connectivity works"
else
    print_warning "HTTP connectivity failed"
fi

echo "Testing HTTPS connectivity..."
if curl -s -o /dev/null -w "%{http_code}" https://$DOMAIN/health --connect-timeout 10 | grep -q "200"; then
    print_success "HTTPS connectivity works"
else
    print_warning "HTTPS connectivity failed"
fi

# Recent events
print_header "Recent Events"
echo "Recent events in namespace $NAMESPACE:"
kubectl get events -n $NAMESPACE --sort-by='.lastTimestamp' | tail -10

echo "Recent cert-manager events:"
kubectl get events -n cert-manager --sort-by='.lastTimestamp' | tail -5

# Summary and recommendations
print_header "Summary and Recommendations"
echo
echo "Configuration Summary:"
echo "- Domain: $DOMAIN"
echo "- Static IP: $STATIC_IP_NAME"
echo "- Ingress: $INGRESS_NAME"
echo "- Certificate: $CERT_NAME (GKE Managed)"
echo "- Namespace: $NAMESPACE"
echo
echo "Common troubleshooting steps:"
echo "1. Ensure DNS points to the correct IP address"
echo "2. Wait for certificate provisioning (can take 10-60 minutes)"
echo "3. Check that all pods are running and healthy"
echo "4. Verify ingress controller is properly configured"
echo "5. Check GCP firewall rules allow traffic on ports 80/443"
echo
echo "Useful commands:"
echo "- Watch ingress: kubectl get ingress $INGRESS_NAME -n $NAMESPACE -w"
echo "- Watch certificate: kubectl get managedcertificate $CERT_NAME -n $NAMESPACE -w"
echo "- Check ingress logs: kubectl logs -n kube-system -l k8s-app=glbc"
echo "- Test endpoint: curl -v https://$DOMAIN/health"