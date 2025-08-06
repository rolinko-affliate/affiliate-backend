#!/bin/bash
set -e

# Default environment is prod
ENV=${1:-prod}

# Validate environment
if [[ "$ENV" != "prod" ]]; then
  echo "Error: Only 'prod' environment is available"
  echo "Usage: $0 [prod]"
  exit 1
fi

echo "Creating secrets for $ENV environment..."

# Create namespace if it doesn't exist
kubectl create namespace affiliate-backend --dry-run=client -o yaml | kubectl apply -f -

# Database credentials secret
echo "Creating database credentials secret..."
echo "Please provide the following database credentials:"

read -p "Database User: " DB_USER
read -s -p "Database Password: " DB_PASSWORD
echo
read -p "Database Name: " DB_NAME
read -p "Cloud SQL Connection Name (project:region:instance): " CONNECTION_NAME

kubectl create secret generic saas-bff-db-credentials \
  --from-literal=db_user="$DB_USER" \
  --from-literal=db_password="$DB_PASSWORD" \
  --from-literal=db_name="$DB_NAME" \
  --from-literal=connection_name="$CONNECTION_NAME" \
  --namespace=affiliate-backend \
  --dry-run=client -o yaml | kubectl apply -f -

# Application secrets
echo "Creating application secrets..."
echo "Please provide the following application secrets:"

read -p "Supabase URL: " SUPABASE_URL
read -s -p "Supabase Anon Key: " SUPABASE_ANON_KEY
echo
read -s -p "Supabase Service Role Key: " SUPABASE_SERVICE_ROLE_KEY
echo
read -s -p "JWT Secret: " JWT_SECRET
echo

kubectl create secret generic saas-bff-secrets \
  --from-literal=supabase_url="$SUPABASE_URL" \
  --from-literal=supabase_anon_key="$SUPABASE_ANON_KEY" \
  --from-literal=supabase_service_role_key="$SUPABASE_SERVICE_ROLE_KEY" \
  --from-literal=jwt_secret="$JWT_SECRET" \
  --namespace=affiliate-backend \
  --dry-run=client -o yaml | kubectl apply -f -

echo "Secrets created successfully for $ENV environment!"
echo "You can now deploy the application using: ./deploy.sh $ENV"