#!/bin/bash
set -e

echo "=== Cloud SQL Proxy Permission Fix Script ==="
echo

# Variables
PROJECT_ID="jinko-test"
NAMESPACE="saas-bff"
CLUSTER_NAME=${CLUSTER_NAME:-"your-cluster-name"}
REGION=${REGION:-"europe-west1"}

echo "Project: $PROJECT_ID"
echo "Namespace: $NAMESPACE"
echo "Cluster: $CLUSTER_NAME"
echo "Region: $REGION"
echo

# Step 1: Check current setup
echo "1. Checking current service account setup..."
SA_ANNOTATION=$(kubectl get serviceaccount default -n $NAMESPACE -o jsonpath='{.metadata.annotations.iam\.gke\.io/gcp-service-account}' 2>/dev/null || echo "")

if [ -z "$SA_ANNOTATION" ]; then
    echo "❌ No Workload Identity annotation found on service account"
    echo "You need to set up Workload Identity. Here's how:"
    echo
    echo "# Create a GCP service account"
    echo "gcloud iam service-accounts create cloudsql-proxy-sa --project=$PROJECT_ID"
    echo
    echo "# Grant Cloud SQL permissions"
    echo "gcloud projects add-iam-policy-binding $PROJECT_ID \\"
    echo "  --member='serviceAccount:cloudsql-proxy-sa@$PROJECT_ID.iam.gserviceaccount.com' \\"
    echo "  --role='roles/cloudsql.client'"
    echo
    echo "# Enable Workload Identity binding"
    echo "gcloud iam service-accounts add-iam-policy-binding \\"
    echo "  cloudsql-proxy-sa@$PROJECT_ID.iam.gserviceaccount.com \\"
    echo "  --role='roles/iam.workloadIdentityUser' \\"
    echo "  --member='serviceAccount:$PROJECT_ID.svc.id.goog[$NAMESPACE/default]'"
    echo
    echo "# Annotate the Kubernetes service account"
    echo "kubectl annotate serviceaccount default -n $NAMESPACE \\"
    echo "  iam.gke.io/gcp-service-account=cloudsql-proxy-sa@$PROJECT_ID.iam.gserviceaccount.com"
    echo
else
    echo "✅ Workload Identity annotation found: $SA_ANNOTATION"
    
    # Check if the GCP service account has the right permissions
    echo "2. Checking IAM permissions for: $SA_ANNOTATION"
    
    # Check if service account has cloudsql.client role
    POLICY_CHECK=$(gcloud projects get-iam-policy $PROJECT_ID --flatten='bindings[].members' --filter="bindings.members:serviceAccount:$SA_ANNOTATION AND bindings.role:roles/cloudsql.client" --format="value(bindings.role)" 2>/dev/null || echo "")
    
    if [ -z "$POLICY_CHECK" ]; then
        echo "❌ Service account missing Cloud SQL client permissions"
        echo "Adding permissions..."
        gcloud projects add-iam-policy-binding $PROJECT_ID \
          --member="serviceAccount:$SA_ANNOTATION" \
          --role="roles/cloudsql.client"
        echo "✅ Permissions added"
    else
        echo "✅ Service account has Cloud SQL client permissions"
    fi
fi

# Step 2: Check Cloud SQL instance
echo
echo "3. Checking Cloud SQL instance..."
INSTANCE_CHECK=$(gcloud sql instances describe jinko-test-shared-postgres --project=$PROJECT_ID --format="value(name)" 2>/dev/null || echo "")

if [ -z "$INSTANCE_CHECK" ]; then
    echo "❌ Cloud SQL instance 'jinko-test-shared-postgres' not found or not accessible"
    echo "Please verify:"
    echo "- Instance name is correct"
    echo "- You have access to the project"
    echo "- Instance is in the correct project"
else
    echo "✅ Cloud SQL instance found: $INSTANCE_CHECK"
fi

# Step 3: Check secrets
echo
echo "4. Checking Kubernetes secrets..."
SECRET_CHECK=$(kubectl get secret saas-bff-db-credentials -n $NAMESPACE -o name 2>/dev/null || echo "")

if [ -z "$SECRET_CHECK" ]; then
    echo "❌ Secret 'saas-bff-db-credentials' not found in namespace $NAMESPACE"
    echo "Create the secret using:"
    echo "./create-app-secrets.sh"
else
    echo "✅ Database credentials secret found"
    CONNECTION_NAME=$(kubectl get secret saas-bff-db-credentials -n $NAMESPACE -o jsonpath='{.data.connection_name}' | base64 -d 2>/dev/null || echo "")
    echo "Connection name in secret: $CONNECTION_NAME"
fi

echo
echo "=== Summary ==="
echo "After making any changes above, restart the deployment:"
echo "kubectl rollout restart deployment/prod-affiliate-backend -n $NAMESPACE"
echo
echo "Monitor the logs:"
echo "kubectl logs -f deployment/prod-affiliate-backend -n $NAMESPACE -c cloud-sql-proxy"