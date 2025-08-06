#!/bin/bash

# Quick Ingress and Certificate Status Check
# Usage: ./ingress-quick-check.sh

NAMESPACE="saas-bff"
DOMAIN="api.affiliate.rolinko.com"

echo "🔍 Quick Ingress & Certificate Status Check"
echo "============================================"

# Check ingress IP
echo "📡 Ingress Status:"
INGRESS_IP=$(kubectl get ingress prod-affiliate-backend -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "No IP")
echo "   IP Address: $INGRESS_IP"

# Check certificate status
echo "🔒 Certificate Status:"
CERT_STATUS=$(kubectl get managedcertificate affiliate-backend-ssl-cert -n $NAMESPACE -o jsonpath='{.status.certificateStatus}' 2>/dev/null || echo "Not found")
echo "   Status: $CERT_STATUS"

# Check DNS resolution
echo "🌐 DNS Resolution:"
DNS_IP=$(nslookup $DOMAIN 2>/dev/null | grep "Address:" | tail -1 | awk '{print $2}' || echo "Failed")
echo "   $DOMAIN → $DNS_IP"

# Check if DNS matches ingress
if [ "$INGRESS_IP" != "No IP" ] && [ "$DNS_IP" != "Failed" ]; then
    if [ "$INGRESS_IP" = "$DNS_IP" ]; then
        echo "   ✅ DNS matches ingress IP"
    else
        echo "   ⚠️  DNS does not match ingress IP"
    fi
fi

# Check pod status
echo "🚀 Pod Status:"
READY_PODS=$(kubectl get pods -n $NAMESPACE -l app=affiliate-backend --no-headers 2>/dev/null | grep "Running" | grep "1/1" | wc -l)
TOTAL_PODS=$(kubectl get pods -n $NAMESPACE -l app=affiliate-backend --no-headers 2>/dev/null | wc -l)
echo "   Ready: $READY_PODS/$TOTAL_PODS"

# Quick connectivity test
echo "🌍 Connectivity Test:"
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://$DOMAIN/health --connect-timeout 5 2>/dev/null || echo "Failed")
HTTPS_STATUS=$(curl -s -o /dev/null -w "%{http_code}" https://$DOMAIN/health --connect-timeout 5 2>/dev/null || echo "Failed")
echo "   HTTP:  $HTTP_STATUS"
echo "   HTTPS: $HTTPS_STATUS"

echo ""
echo "💡 For detailed analysis, run: ./k8s/scripts/check-ingress-and-certs.sh"