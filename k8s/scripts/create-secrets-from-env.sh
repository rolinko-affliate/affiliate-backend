#!/bin/bash
set -e

# Script to create Kubernetes secrets from .env file and additional parameters
# Usage: ./create-secrets-from-env.sh [environment] [--dry-run|--show-yaml]

# Default environment is prod
ENV=${1:-prod}
DRY_RUN=""
SHOW_YAML=""

# Check for flags
if [[ "$2" == "--dry-run" ]] || [[ "$1" == "--dry-run" ]]; then
    DRY_RUN="--dry-run=client"
    echo "🔍 Running in dry-run mode - no actual changes will be made"
elif [[ "$2" == "--show-yaml" ]] || [[ "$1" == "--show-yaml" ]]; then
    SHOW_YAML="true"
    echo "📄 Showing YAML output only - no kubectl commands will be executed"
fi

# Validate environment
if [[ "$ENV" != "prod" ]] && [[ "$ENV" != "--dry-run" ]] && [[ "$ENV" != "--show-yaml" ]]; then
  echo "❌ Error: Only 'prod' environment is available"
  echo "Usage: $0 [prod] [--dry-run|--show-yaml]"
  exit 1
fi

# Handle case where first argument is a flag
if [[ "$ENV" == "--dry-run" ]]; then
    ENV="prod"
    DRY_RUN="--dry-run=client"
    echo "🔍 Running in dry-run mode - no actual changes will be made"
elif [[ "$ENV" == "--show-yaml" ]]; then
    ENV="prod"
    SHOW_YAML="true"
    echo "📄 Showing YAML output only - no kubectl commands will be executed"
fi

echo "🚀 Creating secrets for $ENV environment..."

# Check if .env file exists
ENV_FILE=".env"
if [[ ! -f "$ENV_FILE" ]]; then
    echo "❌ Error: .env file not found at $ENV_FILE"
    exit 1
fi

echo "📄 Reading configuration from $ENV_FILE"

# Source the .env file to load variables
set -a  # automatically export all variables
source "$ENV_FILE"
set +a  # stop automatically exporting

# Create namespace if it doesn't exist
echo "📦 Creating namespace affiliate-backend..."
kubectl create namespace affiliate-backend --dry-run=client -o yaml | kubectl apply -f - $DRY_RUN

# Database credentials secret
echo "🗄️  Creating database credentials secret..."

# For production, these should be provided as environment variables or parameters
# For now, using default values that can be overridden
DB_USER=${DB_USER:-"postgres"}
DB_PASSWORD=${DB_PASSWORD:-"postgres"}
DB_NAME=${DB_NAME:-"affiliate_platform"}
CONNECTION_NAME=${CONNECTION_NAME:-"jinko-test:europe-west1:affiliate-db"}

echo "   Database User: $DB_USER"
echo "   Database Name: $DB_NAME"
echo "   Connection Name: $CONNECTION_NAME"

kubectl create secret generic saas-bff-db-credentials \
  --from-literal=db_user="$DB_USER" \
  --from-literal=db_password="$DB_PASSWORD" \
  --from-literal=db_name="$DB_NAME" \
  --from-literal=connection_name="$CONNECTION_NAME" \
  --namespace=affiliate-backend \
  --dry-run=client -o yaml | kubectl apply -f - $DRY_RUN

# Application secrets from .env file
echo "🔐 Creating application secrets from .env file..."

# Map .env variables to Kubernetes secret keys
JWT_SECRET_VALUE=${SUPABASE_JWT_SECRET}
ENCRYPTION_KEY_VALUE=${ENCRYPTION_KEY}

# These need to be provided as environment variables since they're not in .env
SUPABASE_URL_VALUE=${SUPABASE_URL:-"https://your-project.supabase.co"}
SUPABASE_ANON_KEY_VALUE=${SUPABASE_ANON_KEY:-"your-anon-key"}
SUPABASE_SERVICE_ROLE_KEY_VALUE=${SUPABASE_SERVICE_ROLE_KEY:-"your-service-role-key"}

echo "   JWT Secret: ${JWT_SECRET_VALUE:0:20}... (truncated)"
echo "   Encryption Key: ${ENCRYPTION_KEY_VALUE:0:20}... (truncated)"
echo "   Supabase URL: $SUPABASE_URL_VALUE"
echo "   Supabase Anon Key: ${SUPABASE_ANON_KEY_VALUE:0:20}... (truncated)"
echo "   Supabase Service Role Key: ${SUPABASE_SERVICE_ROLE_KEY_VALUE:0:20}... (truncated)"

kubectl create secret generic saas-bff-secrets \
  --from-literal=supabase_url="$SUPABASE_URL_VALUE" \
  --from-literal=supabase_anon_key="$SUPABASE_ANON_KEY_VALUE" \
  --from-literal=supabase_service_role_key="$SUPABASE_SERVICE_ROLE_KEY_VALUE" \
  --from-literal=jwt_secret="$JWT_SECRET_VALUE" \
  --from-literal=encryption_key="$ENCRYPTION_KEY_VALUE" \
  --namespace=affiliate-backend \
  --dry-run=client -o yaml | kubectl apply -f - $DRY_RUN

if [[ -z "$DRY_RUN" ]]; then
    echo "✅ Secrets created successfully for $ENV environment!"
    echo ""
    echo "📋 Summary of created secrets:"
    echo "   • saas-bff-db-credentials (database connection info)"
    echo "   • saas-bff-secrets (application secrets)"
    echo ""
    echo "🚀 You can now deploy the application using: ./deploy.sh $ENV"
    echo ""
    echo "🔍 To verify secrets were created:"
    echo "   kubectl get secrets -n affiliate-backend"
    echo "   kubectl describe secret saas-bff-secrets -n affiliate-backend"
else
    echo "✅ Dry-run completed successfully!"
    echo "   Remove --dry-run flag to actually create the secrets"
fi

echo ""
echo "⚠️  IMPORTANT NOTES:"
echo "   • For production, set these environment variables before running:"
echo "     export SUPABASE_URL='https://your-project.supabase.co'"
echo "     export SUPABASE_ANON_KEY='your-actual-anon-key'"
echo "     export SUPABASE_SERVICE_ROLE_KEY='your-actual-service-role-key'"
echo "     export DB_PASSWORD='your-actual-db-password'"
echo "     export CONNECTION_NAME='your-project:region:instance'"
echo "   • Current values are using defaults/placeholders"
echo "   • Secrets contain sensitive data - handle with care"