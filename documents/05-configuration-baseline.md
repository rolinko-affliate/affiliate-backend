# Configuration Baseline (CMDB / IaC): Affiliate Backend Platform

**Document Version**: v1.1  
**Owner**: DevOps Engineer  
**Last Updated**: 2025-08-15  
**Next Review**: 2026-02-15

---

## 1. Executive Summary

This document serves as the **Configuration Management Database (CMDB)** and **Infrastructure as Code (IaC)** baseline for the Affiliate Backend Platform. It provides the authoritative source of truth for all infrastructure components, their configurations, dependencies, and change management processes.

### Key Principles
- **Infrastructure as Code**: All infrastructure defined in version-controlled code
- **Immutable Infrastructure**: Replace rather than modify infrastructure components
- **Configuration Drift Prevention**: Automated detection and remediation of configuration changes
- **Audit Trail**: Complete history of all infrastructure changes

## 2. Environment Inventory

### 2.1 Environment Overview

| Environment | Purpose | Region | Kubernetes Cluster | Database | Status |
|-------------|---------|--------|-------------------|----------|--------|
| **Production** | Live customer service | us-central1-a | affiliate-prod-gke | affiliate-prod-db | Active |
| **Staging** | Pre-production testing | us-west1-a | affiliate-staging-gke | affiliate-staging-db | Active |
| **Development** | Local development | N/A | Docker Compose | Local PostgreSQL | Active |

### 2.2 Production Environment Details

#### Google Cloud Platform Resources
```yaml
Project ID: affiliate-platform-prod
Region: us-central1
Zone: us-central1-a

Compute Resources:
  - GKE Cluster: affiliate-prod-gke
    - Node Pool: default-pool
    - Machine Type: e2-standard-4
    - Nodes: 3-10 (auto-scaling)
    - Disk: 100GB SSD per node
  
Database Resources:
  - Cloud SQL: affiliate-prod-db
    - Instance Type: db-custom-4-16384
    - Storage: 500GB SSD
    - Backup: Automated daily
    - High Availability: Enabled
  
  - Memorystore Redis: affiliate-prod-redis
    - Memory: 5GB
    - Version: Redis 6.2
    - High Availability: Enabled

Storage Resources:
  - Cloud Storage: affiliate-prod-backups
    - Storage Class: Standard
    - Lifecycle: 90 days retention
  
  - Cloud Storage: affiliate-prod-assets
    - Storage Class: Standard
    - CDN: Enabled

Network Resources:
  - VPC: affiliate-prod-vpc
  - Subnets: 
    - gke-subnet: 10.0.0.0/16
    - db-subnet: 10.1.0.0/24
  - Load Balancer: affiliate-prod-lb
  - Cloud NAT: affiliate-prod-nat
```

### 2.3 Staging Environment Details

#### Google Cloud Platform Resources
```yaml
Project ID: affiliate-platform-staging
Region: us-west1
Zone: us-west1-a

Compute Resources:
  - GKE Cluster: affiliate-staging-gke
    - Node Pool: default-pool
    - Machine Type: e2-standard-2
    - Nodes: 2 (fixed)
    - Disk: 50GB SSD per node

Database Resources:
  - Cloud SQL: affiliate-staging-db
    - Instance Type: db-custom-2-8192
    - Storage: 100GB SSD
    - Backup: Automated daily
    - High Availability: Disabled
  
  - Memorystore Redis: affiliate-staging-redis
    - Memory: 1GB
    - Version: Redis 6.2
    - High Availability: Disabled

Storage Resources:
  - Cloud Storage: affiliate-staging-backups
    - Storage Class: Standard
    - Lifecycle: 30 days retention

Network Resources:
  - VPC: affiliate-staging-vpc
  - Subnets:
    - gke-subnet: 10.10.0.0/16
    - db-subnet: 10.11.0.0/24
  - Load Balancer: affiliate-staging-lb
```

## 3. Infrastructure as Code Repository Structure

### 3.1 Repository Organization

```
infrastructure/
├── terraform/
│   ├── environments/
│   │   ├── production/
│   │   │   ├── main.tf
│   │   │   ├── variables.tf
│   │   │   ├── outputs.tf
│   │   │   └── terraform.tfvars
│   │   ├── staging/
│   │   │   ├── main.tf
│   │   │   ├── variables.tf
│   │   │   ├── outputs.tf
│   │   │   └── terraform.tfvars
│   │   └── modules/
│   │       ├── gke-cluster/
│   │       ├── cloud-sql/
│   │       ├── redis/
│   │       ├── networking/
│   │       └── monitoring/
│   └── shared/
│       ├── backend.tf
│       ├── providers.tf
│       └── versions.tf
├── kubernetes/
│   ├── base/
│   │   ├── namespace.yaml
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── ingress.yaml
│   │   └── kustomization.yaml
│   ├── overlays/
│   │   ├── production/
│   │   │   ├── kustomization.yaml
│   │   │   ├── replica-patch.yaml
│   │   │   ├── resource-patch.yaml
│   │   │   └── config-patch.yaml
│   │   └── staging/
│   │       ├── kustomization.yaml
│   │       ├── replica-patch.yaml
│   │       └── config-patch.yaml
│   └── monitoring/
│       ├── prometheus/
│       ├── grafana/
│       └── alertmanager/
├── scripts/
│   ├── deploy.sh
│   ├── rollback.sh
│   └── validate.sh
└── docs/
    ├── deployment-guide.md
    └── troubleshooting.md
```

### 3.2 Repository Links

| Repository | URL | Purpose | Access Level |
|------------|-----|---------|--------------|
| **Infrastructure** | `https://github.com/company/affiliate-infrastructure` | Terraform & K8s configs | DevOps Team |
| **Application** | `https://github.com/company/affiliate-backend` | Application source code | Development Team |
| **Monitoring** | `https://github.com/company/affiliate-monitoring` | Monitoring configurations | SRE Team |
| **Documentation** | `https://github.com/company/affiliate-docs` | Service documentation | All Teams |

## 4. Configuration Management

### 4.1 Terraform Configuration

#### Production Environment (`terraform/environments/production/main.tf`)
```hcl
terraform {
  required_version = ">= 1.5"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.84"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
  }
  
  backend "gcs" {
    bucket = "affiliate-terraform-state-prod"
    prefix = "terraform/state"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# GKE Cluster
module "gke_cluster" {
  source = "../../modules/gke-cluster"
  
  project_id   = var.project_id
  region       = var.region
  cluster_name = "affiliate-prod-gke"
  
  node_pools = [
    {
      name         = "default-pool"
      machine_type = "e2-standard-4"
      min_count    = 3
      max_count    = 10
      disk_size_gb = 100
      disk_type    = "pd-ssd"
    }
  ]
  
  network_config = {
    vpc_name     = "affiliate-prod-vpc"
    subnet_name  = "gke-subnet"
    subnet_cidr  = "10.0.0.0/16"
  }
}

# Cloud SQL Database
module "cloud_sql" {
  source = "../../modules/cloud-sql"
  
  project_id      = var.project_id
  region          = var.region
  instance_name   = "affiliate-prod-db"
  database_name   = "affiliate_platform"
  
  instance_config = {
    tier                = "db-custom-4-16384"
    disk_size          = 500
    disk_type          = "PD_SSD"
    availability_type  = "REGIONAL"
    backup_enabled     = true
    maintenance_window = {
      day  = 7  # Sunday
      hour = 3  # 3 AM
    }
  }
  
  database_flags = [
    {
      name  = "max_connections"
      value = "200"
    },
    {
      name  = "shared_preload_libraries"
      value = "pg_stat_statements"
    }
  ]
}

# Redis Instance
module "redis" {
  source = "../../modules/redis"
  
  project_id    = var.project_id
  region        = var.region
  instance_name = "affiliate-prod-redis"
  
  config = {
    memory_size_gb     = 5
    redis_version      = "REDIS_6_X"
    availability_type  = "STANDARD_HA"
    auth_enabled      = true
  }
}

# Monitoring
module "monitoring" {
  source = "../../modules/monitoring"
  
  project_id     = var.project_id
  cluster_name   = module.gke_cluster.cluster_name
  notification_channels = var.notification_channels
}
```

#### Variables (`terraform/environments/production/variables.tf`)
```hcl
variable "project_id" {
  description = "GCP Project ID"
  type        = string
  default     = "affiliate-platform-prod"
}

variable "region" {
  description = "GCP Region"
  type        = string
  default     = "us-central1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "notification_channels" {
  description = "List of notification channels for alerts"
  type        = list(string)
  default     = [
    "projects/affiliate-platform-prod/notificationChannels/1234567890",
    "projects/affiliate-platform-prod/notificationChannels/0987654321"
  ]
}
```

### 4.2 Kubernetes Configuration

#### Base Deployment (`kubernetes/base/deployment.yaml`)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: affiliate-api
  labels:
    app: affiliate-api
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: affiliate-api
  template:
    metadata:
      labels:
        app: affiliate-api
        version: v1
    spec:
      containers:
      - name: api
        image: gcr.io/affiliate-platform-prod/affiliate-api:latest
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: PORT
          value: "8080"
        - name: ENVIRONMENT
          value: "production"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: database-secret
              key: url
        - name: SUPABASE_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: auth-secret
              key: jwt-secret
        - name: ENCRYPTION_KEY
          valueFrom:
            secretKeyRef:
              name: encryption-secret
              key: key
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        securityContext:
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          runAsUser: 1000
          capabilities:
            drop:
            - ALL
```

#### Production Overlay (`kubernetes/overlays/production/kustomization.yaml`)
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: affiliate-prod

resources:
- ../../base

patchesStrategicMerge:
- replica-patch.yaml
- resource-patch.yaml

configMapGenerator:
- name: app-config
  literals:
  - ENVIRONMENT=production
  - DEBUG_MODE=false
  - MOCK_MODE=false

secretGenerator:
- name: database-secret
  literals:
  - url=postgres://user:pass@affiliate-prod-db:5432/affiliate_platform?sslmode=require

images:
- name: gcr.io/affiliate-platform-prod/affiliate-api
  newTag: v1.2.3
```

## 5. Version Control and Approval Workflow

### 5.1 Infrastructure Change Process

```mermaid
graph TD
    DEV[Developer Creates PR] --> REVIEW[Code Review]
    REVIEW --> PLAN[Terraform Plan]
    PLAN --> APPROVE{Approval Required?}
    APPROVE -->|Minor Change| AUTO[Auto-Apply]
    APPROVE -->|Major Change| MANUAL[Manual Approval]
    MANUAL --> APPLY[Terraform Apply]
    AUTO --> APPLY
    APPLY --> VERIFY[Verify Deployment]
    VERIFY --> MONITOR[Monitor Changes]
    MONITOR --> COMPLETE[Change Complete]
```

### 5.2 Change Classification

| Change Type | Examples | Approval Required | Auto-Apply | Rollback Plan |
|-------------|----------|------------------|------------|---------------|
| **Minor** | Resource scaling, config updates | Tech Lead | Yes | Automated |
| **Major** | New services, network changes | Service Manager + SRE Lead | No | Manual |
| **Critical** | Security updates, DB schema | All stakeholders | No | Full DR test |
| **Emergency** | Hotfixes, security patches | Service Manager | Yes | Immediate rollback |

### 5.3 Approval Matrix

| Resource Type | Development | Staging | Production |
|---------------|-------------|---------|------------|
| **Application Config** | Developer | Tech Lead | Service Manager |
| **Infrastructure** | DevOps Engineer | SRE Lead | Service Manager + SRE Lead |
| **Database Schema** | Developer | Tech Lead | Service Manager + DBA |
| **Security Config** | Security Engineer | Security Officer | CISO |
| **Network Changes** | DevOps Engineer | SRE Lead | Network Admin + SRE Lead |

## 6. Configuration Drift Detection

### 6.1 Automated Drift Detection

#### Terraform Drift Detection
```bash
#!/bin/bash
# scripts/drift-detection.sh

set -e

ENVIRONMENTS=("production" "staging")
SLACK_WEBHOOK_URL="${SLACK_WEBHOOK_URL}"

for env in "${ENVIRONMENTS[@]}"; do
    echo "Checking drift for environment: $env"
    
    cd "terraform/environments/$env"
    
    # Initialize Terraform
    terraform init -backend=true
    
    # Plan and capture output
    if ! terraform plan -detailed-exitcode -out=tfplan > plan_output.txt 2>&1; then
        # Drift detected (exit code 2) or error (exit code 1)
        echo "Configuration drift detected in $env environment"
        
        # Send alert to Slack
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"🚨 Configuration drift detected in $env environment. Check the logs for details.\"}" \
            "$SLACK_WEBHOOK_URL"
        
        # Save drift report
        cp plan_output.txt "../../../reports/drift-$env-$(date +%Y%m%d-%H%M%S).txt"
    else
        echo "No drift detected in $env environment"
    fi
    
    # Cleanup
    rm -f tfplan plan_output.txt
    cd - > /dev/null
done
```

#### Kubernetes Configuration Drift
```bash
#!/bin/bash
# scripts/k8s-drift-detection.sh

set -e

ENVIRONMENTS=("production" "staging")

for env in "${ENVIRONMENTS[@]}"; do
    echo "Checking Kubernetes drift for environment: $env"
    
    # Switch to appropriate context
    kubectl config use-context "affiliate-$env-gke"
    
    # Check for configuration drift
    kubectl diff -k "kubernetes/overlays/$env" > "drift-$env.txt" || true
    
    if [ -s "drift-$env.txt" ]; then
        echo "Kubernetes configuration drift detected in $env"
        # Send alert and save report
        curl -X POST -H 'Content-type: application/json' \
            --data "{\"text\":\"🚨 Kubernetes configuration drift detected in $env environment.\"}" \
            "$SLACK_WEBHOOK_URL"
    else
        echo "No Kubernetes drift detected in $env"
    fi
    
    rm -f "drift-$env.txt"
done
```

### 6.2 Drift Remediation Process

1. **Detection**: Automated daily scans detect configuration drift
2. **Alert**: Immediate notification to DevOps team via Slack/PagerDuty
3. **Assessment**: Team reviews drift to determine if intentional or accidental
4. **Remediation**: 
   - If accidental: Apply IaC to restore desired state
   - If intentional: Update IaC to match current state
5. **Documentation**: Update change log and notify stakeholders

## 7. Backup and Recovery Configuration

### 7.1 Database Backup Strategy

#### Automated Backups
```yaml
# Cloud SQL Backup Configuration
backup_configuration:
  enabled: true
  start_time: "03:00"  # 3 AM UTC
  location: "us-central1"
  point_in_time_recovery_enabled: true
  transaction_log_retention_days: 7
  backup_retention_settings:
    retained_backups: 30
    retention_unit: "COUNT"
```

#### Manual Backup Process
```bash
#!/bin/bash
# scripts/manual-backup.sh

ENVIRONMENT=${1:-production}
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_NAME="manual-backup-$ENVIRONMENT-$TIMESTAMP"

echo "Creating manual backup for $ENVIRONMENT environment"

gcloud sql backups create \
    --instance="affiliate-$ENVIRONMENT-db" \
    --description="Manual backup created on $TIMESTAMP" \
    --project="affiliate-platform-$ENVIRONMENT"

echo "Backup created: $BACKUP_NAME"
```

### 7.2 Application State Backup

#### Kubernetes State Backup
```bash
#!/bin/bash
# scripts/k8s-backup.sh

ENVIRONMENT=${1:-production}
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
BACKUP_DIR="backups/k8s-$ENVIRONMENT-$TIMESTAMP"

mkdir -p "$BACKUP_DIR"

# Backup all resources
kubectl get all --all-namespaces -o yaml > "$BACKUP_DIR/all-resources.yaml"
kubectl get configmaps --all-namespaces -o yaml > "$BACKUP_DIR/configmaps.yaml"
kubectl get secrets --all-namespaces -o yaml > "$BACKUP_DIR/secrets.yaml"
kubectl get pv,pvc --all-namespaces -o yaml > "$BACKUP_DIR/volumes.yaml"

# Compress and upload to Cloud Storage
tar -czf "$BACKUP_DIR.tar.gz" "$BACKUP_DIR"
gsutil cp "$BACKUP_DIR.tar.gz" "gs://affiliate-$ENVIRONMENT-backups/kubernetes/"

echo "Kubernetes backup completed: $BACKUP_DIR.tar.gz"
```

## 8. Monitoring and Alerting Configuration

### 8.1 Infrastructure Monitoring

#### Prometheus Configuration (`monitoring/prometheus/prometheus.yml`)
```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "rules/*.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

scrape_configs:
  - job_name: 'kubernetes-apiservers'
    kubernetes_sd_configs:
    - role: endpoints
    scheme: https
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
    - source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
      action: keep
      regex: default;kubernetes;https

  - job_name: 'kubernetes-nodes'
    kubernetes_sd_configs:
    - role: node
    scheme: https
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
    - action: labelmap
      regex: __meta_kubernetes_node_label_(.+)

  - job_name: 'kubernetes-pods'
    kubernetes_sd_configs:
    - role: pod
    relabel_configs:
    - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
      action: keep
      regex: true
    - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
      action: replace
      target_label: __metrics_path__
      regex: (.+)
```

#### Alert Rules (`monitoring/prometheus/rules/infrastructure.yml`)
```yaml
groups:
- name: infrastructure
  rules:
  - alert: HighCPUUsage
    expr: 100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High CPU usage detected"
      description: "CPU usage is above 80% for more than 5 minutes"

  - alert: HighMemoryUsage
    expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 85
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "High memory usage detected"
      description: "Memory usage is above 85% for more than 5 minutes"

  - alert: DatabaseConnectionsHigh
    expr: pg_stat_database_numbackends / pg_settings_max_connections * 100 > 80
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "Database connections high"
      description: "Database connection usage is above 80%"
```

### 8.2 Application Monitoring

#### Grafana Dashboard Configuration
```json
{
  "dashboard": {
    "title": "Affiliate Backend Platform - Production",
    "panels": [
      {
        "title": "API Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{status}}"
          }
        ]
      },
      {
        "title": "API Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      },
      {
        "title": "Database Query Performance",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(pg_stat_database_tup_returned[5m])",
            "legendFormat": "Rows returned/sec"
          }
        ]
      }
    ]
  }
}
```

## 9. Security Configuration

### 9.1 Network Security

#### VPC Configuration
```hcl
# Network security configuration
resource "google_compute_network" "vpc" {
  name                    = "affiliate-${var.environment}-vpc"
  auto_create_subnetworks = false
  routing_mode           = "REGIONAL"
}

resource "google_compute_subnetwork" "gke_subnet" {
  name          = "gke-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.region
  network       = google_compute_network.vpc.id
  
  secondary_ip_range {
    range_name    = "pods"
    ip_cidr_range = "10.1.0.0/16"
  }
  
  secondary_ip_range {
    range_name    = "services"
    ip_cidr_range = "10.2.0.0/16"
  }
}

resource "google_compute_firewall" "allow_internal" {
  name    = "allow-internal"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "icmp"
  }

  source_ranges = ["10.0.0.0/8"]
}
```

### 9.2 Secret Management

#### Kubernetes Secrets
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: database-secret
  namespace: affiliate-prod
type: Opaque
data:
  url: <base64-encoded-database-url>
  username: <base64-encoded-username>
  password: <base64-encoded-password>

---
apiVersion: v1
kind: Secret
metadata:
  name: auth-secret
  namespace: affiliate-prod
type: Opaque
data:
  jwt-secret: <base64-encoded-jwt-secret>

---
apiVersion: v1
kind: Secret
metadata:
  name: encryption-secret
  namespace: affiliate-prod
type: Opaque
data:
  key: <base64-encoded-encryption-key>
```

## 10. Disaster Recovery Configuration

### 10.1 Multi-Region Setup

#### Production DR Configuration
```hcl
# Disaster Recovery in us-west1
module "dr_gke_cluster" {
  source = "../../modules/gke-cluster"
  
  project_id   = var.project_id
  region       = "us-west1"
  cluster_name = "affiliate-prod-dr-gke"
  
  node_pools = [
    {
      name         = "dr-pool"
      machine_type = "e2-standard-2"
      min_count    = 0  # Scaled to 0 normally
      max_count    = 5
      disk_size_gb = 100
      disk_type    = "pd-ssd"
    }
  ]
}

module "dr_cloud_sql" {
  source = "../../modules/cloud-sql"
  
  project_id      = var.project_id
  region          = "us-west1"
  instance_name   = "affiliate-prod-dr-db"
  
  replica_configuration = {
    master_instance_name = module.cloud_sql.instance_name
    replica_type        = "READ_REPLICA"
  }
}
```

### 10.2 Backup Retention Policy

| Backup Type | Retention Period | Storage Location | Recovery Time |
|-------------|------------------|------------------|---------------|
| **Database Daily** | 30 days | Regional | 15 minutes |
| **Database Weekly** | 12 weeks | Multi-regional | 30 minutes |
| **Database Monthly** | 12 months | Multi-regional | 1 hour |
| **Application Config** | 90 days | Git + Cloud Storage | 5 minutes |
| **Infrastructure State** | 90 days | Terraform Cloud | 10 minutes |

---

## Appendix A: Environment Variables Reference

### Production Environment Variables
```bash
# Application Configuration
PORT=8080
ENVIRONMENT=production
DEBUG_MODE=false
MOCK_MODE=false

# Database Configuration
DATABASE_URL=postgres://user:pass@affiliate-prod-db:5432/affiliate_platform?sslmode=require
DATABASE_HOST=affiliate-prod-db
DATABASE_PORT=5432
DATABASE_NAME=affiliate_platform
DATABASE_USER=affiliate_user
DATABASE_SSL_MODE=require

# Authentication
SUPABASE_JWT_SECRET=<production-jwt-secret>

# Encryption
ENCRYPTION_KEY=<production-encryption-key>

# External Services
STRIPE_SECRET_KEY=sk_live_<production-stripe-key>
EVERFLOW_API_KEY=<production-everflow-key>
EVERFLOW_BASE_URL=https://api.everflow.io

# Dashboard Configuration
DASHBOARD_CACHE_TTL=300
DASHBOARD_MAX_CONCURRENT_REQUESTS=50
DASHBOARD_TIMEOUT_SECONDS=30
DASHBOARD_RATE_LIMIT_PER_MINUTE=100

# Monitoring
PROMETHEUS_ENDPOINT=http://prometheus:9090
GRAFANA_ENDPOINT=http://grafana:3000
```

## Appendix B: Resource Sizing Guidelines

### Kubernetes Resource Requests and Limits
```yaml
# Production sizing
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "1Gi"
    cpu: "1000m"

# Staging sizing
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

### Database Sizing
| Environment | Instance Type | vCPUs | Memory | Storage | IOPS |
|-------------|---------------|-------|--------|---------|------|
| **Production** | db-custom-4-16384 | 4 | 16 GB | 500 GB SSD | 3,000 |
| **Staging** | db-custom-2-8192 | 2 | 8 GB | 100 GB SSD | 1,500 |
| **Development** | Local PostgreSQL | 2 | 4 GB | 50 GB | N/A |

## Appendix C: Change Log Template

```markdown
# Infrastructure Change Log

## Change ID: CHG-2025-001
**Date**: 2025-08-05  
**Environment**: Production  
**Type**: Major  
**Requestor**: DevOps Engineer  
**Approver**: Service Manager  

### Description
Upgrade GKE cluster from version 1.27 to 1.28

### Components Affected
- GKE Cluster: affiliate-prod-gke
- Node Pool: default-pool
- Kubernetes version: 1.27.x → 1.28.x

### Risk Assessment
- **Risk Level**: Medium
- **Impact**: Temporary pod restarts during node upgrades
- **Mitigation**: Rolling upgrade with pod disruption budgets

### Rollback Plan
1. Revert cluster version using Terraform
2. Drain and recreate nodes if necessary
3. Verify application functionality

### Verification Steps
- [ ] Cluster upgrade completed successfully
- [ ] All pods running and healthy
- [ ] Application endpoints responding
- [ ] Monitoring alerts cleared
```

---

**Document Classification**: Infrastructure Configuration  
**Access Level**: DevOps Team, SRE Team  
**Backup Location**: Git repository + Cloud Storage  
**Related Documents**: System Architecture, Security Guide, Runbook