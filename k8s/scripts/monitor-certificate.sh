#!/bin/bash

# Monitor certificate provisioning progress
NAMESPACE="saas-bff"
CERT_NAME="affiliate-backend-ssl-cert"
DOMAIN="api.affiliate.rolinko.com"

echo "🔍 Monitoring Certificate Provisioning"
echo "======================================"
echo "Certificate: $CERT_NAME"
echo "Domain: $DOMAIN"
echo "Namespace: $NAMESPACE"
echo ""
echo "Press Ctrl+C to stop monitoring"
echo ""

while true; do
    # Clear screen and show timestamp
    clear
    echo "🔍 Certificate Monitoring - $(date)"
    echo "=================================="
    
    # Check if certificate exists
    if kubectl get managedcertificate $CERT_NAME -n $NAMESPACE &>/dev/null; then
        # Get certificate status
        CERT_STATUS=$(kubectl get managedcertificate $CERT_NAME -n $NAMESPACE -o jsonpath='{.status.certificateStatus}' 2>/dev/null || echo "Unknown")
        
        echo "📋 Certificate Status: $CERT_STATUS"
        
        case $CERT_STATUS in
            "Active")
                echo "✅ Certificate is ACTIVE and ready!"
                echo ""
                echo "Testing HTTPS connectivity..."
                HTTPS_STATUS=$(curl -s -o /dev/null -w "%{http_code}" https://$DOMAIN/health --connect-timeout 10 2>/dev/null || echo "Failed")
                echo "HTTPS Status: $HTTPS_STATUS"
                
                if [ "$HTTPS_STATUS" = "200" ]; then
                    echo "🎉 HTTPS is working! Setup complete."
                    break
                fi
                ;;
            "Provisioning")
                echo "⏳ Certificate is being provisioned..."
                echo "   This typically takes 10-60 minutes"
                ;;
            "Failed")
                echo "❌ Certificate provisioning failed!"
                echo ""
                echo "Certificate details:"
                kubectl describe managedcertificate $CERT_NAME -n $NAMESPACE
                break
                ;;
            *)
                echo "⚠️  Unknown status: $CERT_STATUS"
                ;;
        esac
        
        echo ""
        echo "📊 Certificate Details:"
        kubectl get managedcertificate $CERT_NAME -n $NAMESPACE -o wide
        
        echo ""
        echo "🌐 DNS Check:"
        DNS_IP=$(nslookup $DOMAIN 2>/dev/null | grep "Address:" | tail -1 | awk '{print $2}' || echo "Failed")
        INGRESS_IP=$(kubectl get ingress prod-affiliate-backend -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "No IP")
        echo "   DNS: $DOMAIN → $DNS_IP"
        echo "   Ingress IP: $INGRESS_IP"
        
        if [ "$DNS_IP" = "$INGRESS_IP" ]; then
            echo "   ✅ DNS matches ingress IP"
        else
            echo "   ⚠️  DNS does not match ingress IP"
        fi
        
    else
        echo "❌ ManagedCertificate '$CERT_NAME' not found!"
        echo ""
        echo "Create it with:"
        echo "   kubectl apply -f k8s/overlays/prod/managed-certificate.yaml"
        break
    fi
    
    echo ""
    echo "🔄 Refreshing in 30 seconds... (Ctrl+C to stop)"
    sleep 30
done