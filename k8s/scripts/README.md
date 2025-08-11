# Kubernetes Secrets Management Scripts

This directory contains scripts for managing Kubernetes secrets for the affiliate backend application.

## Scripts Overview

### `create-app-secrets.sh`

Creates the `saas-bff-secrets` secret from the `.env` file and environment variables.

**Note**: The `saas-bff-db-credentials` secret is managed by Terraform and should not be created manually.

## Usage

### Basic Usage

```bash
# Create secrets for production environment
./create-app-secrets.sh prod

# Show what would be created without applying
./create-app-secrets.sh --show-yaml

# Dry run (validates against cluster if available)
./create-app-secrets.sh --dry-run
```

### With Environment Variables

For production deployment, set the Supabase credentials as environment variables:

```bash
export SUPABASE_URL="https://your-project.supabase.co"
export SUPABASE_ANON_KEY="your-actual-anon-key"
export SUPABASE_SERVICE_ROLE_KEY="your-actual-service-role-key"

./create-app-secrets.sh prod
```

## Secret Contents

### `saas-bff-secrets`

This secret contains application-level secrets:

| Key | Source | Description |
|-----|--------|-------------|
| `jwt_secret` | `.env` file (`SUPABASE_JWT_SECRET`) | JWT signing secret |
| `encryption_key` | `.env` file (`ENCRYPTION_KEY`) | Application encryption key |
| `supabase_url` | Environment variable | Supabase project URL |
| `supabase_anon_key` | Environment variable | Supabase anonymous key |
| `supabase_service_role_key` | Environment variable | Supabase service role key |

### `saas-bff-db-credentials` (Managed by Terraform)

This secret is automatically created by Terraform and contains:

- `db_user` - Database username
- `db_password` - Database password  
- `db_name` - Database name
- `connection_name` - Cloud SQL connection string

## Examples

### Development/Testing

```bash
# Show the YAML that would be generated
./create-app-secrets.sh --show-yaml
```

### Production Deployment

```bash
# Set production Supabase credentials
export SUPABASE_URL="https://your-prod-project.supabase.co"
export SUPABASE_ANON_KEY="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
export SUPABASE_SERVICE_ROLE_KEY="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Create the secret
./create-app-secrets.sh prod
```

### Verification

After creating secrets, verify they exist:

```bash
# List all secrets in the namespace
kubectl get secrets -n affiliate-backend

# Describe the application secrets
kubectl describe secret saas-bff-secrets -n affiliate-backend

# View secret data (base64 encoded)
kubectl get secret saas-bff-secrets -n affiliate-backend -o yaml
```

## Security Notes

- **Never commit real secrets to version control**
- The `.env` file contains development secrets only
- Production secrets should be set via environment variables
- Use `--show-yaml` to review secrets before applying
- Secrets are base64 encoded in Kubernetes but not encrypted at rest by default
- Consider using tools like Sealed Secrets or External Secrets Operator for production

## Troubleshooting

### Cannot connect to Kubernetes cluster

If you see connection errors:

1. Ensure you're connected to the correct Kubernetes cluster
2. Verify your `kubectl` context: `kubectl config current-context`
3. Use `--show-yaml` to generate the YAML without applying

### Missing environment variables

If Supabase credentials are not set, the script will use placeholder values. Set the required environment variables before running in production.

### Permission errors

Ensure your Kubernetes user has permission to:
- Create secrets in the `affiliate-backend` namespace
- Create namespaces (if the namespace doesn't exist)

## Related Files

- `/workspace/.env` - Development environment variables
- `/workspace/k8s/base/deployment.yaml` - Deployment configuration that uses these secrets
- `/workspace/k8s/overlays/prod/` - Production overlay configuration