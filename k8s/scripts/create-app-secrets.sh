#!/bin/bash
set -e

SUPABASE_JWT_SECRET="gDxsm/JerlPJiOObQLtfjViLBQF2ggmJpYCNW+9LPwL2QJksmiYlzRCJCKseCLxJtGysx+awZvoiS0MF0pLjnw=="
ENCRYPTION_KEY="MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDE="
EVERFLOW_API_KEY="${EVERFLOW_API_KEY:-your-everflow-api-key-here}"

# Generate the kubectl command
KUBECTL_CMD="kubectl create secret generic saas-bff-secrets \
  --from-literal=supabase_jwt_secret=\"$SUPABASE_JWT_SECRET\" \
  --from-literal=encryption_key=\"$ENCRYPTION_KEY\" \
  --from-literal=everflow_api_key=\"$EVERFLOW_API_KEY\" \
  --namespace=saas-bff"

eval "$KUBECTL_CMD --dry-run=client -o yaml | kubectl apply -f -"
