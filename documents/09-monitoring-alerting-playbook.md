# Monitoring & Alerting Playbook: Affiliate Backend Platform

**Document Version**: v1.0  
**Owner**: SRE Team  
**Last Updated**: 2025-08-05  
**Next Review**: 2026-02-05

---

## 1. Overview

This playbook defines the monitoring and alerting strategy for the Affiliate Backend Platform. It provides comprehensive guidance on monitoring infrastructure, applications, and business metrics to ensure service reliability and performance while maintaining our 99.9% availability SLA.

### Monitoring Objectives
- **Proactive Issue Detection**: Identify problems before they impact customers
- **SLA Compliance**: Monitor and maintain 99.9% availability target
- **Performance Optimization**: Track and improve system performance
- **Business Intelligence**: Monitor key business metrics and trends

### Key Monitoring Principles
- **Four Golden Signals**: Latency, Traffic, Errors, Saturation
- **SLI/SLO Driven**: Monitor what matters to users
- **Actionable Alerts**: Every alert should require human action
- **Escalation Clarity**: Clear escalation paths for all alerts

## 2. Monitoring Architecture

### 2.1 Monitoring Stack Overview

```mermaid
graph TB
    subgraph "Data Collection"
        APP[Application Metrics<br/>Prometheus Client]
        INFRA[Infrastructure Metrics<br/>Node Exporter]
        K8S[Kubernetes Metrics<br/>kube-state-metrics]
        LOGS[Application Logs<br/>Structured JSON]
        TRACES[Distributed Traces<br/>Jaeger]
    end
    
    subgraph "Storage & Processing"
        PROM[(Prometheus<br/>Time Series DB)]
        ELASTIC[(Elasticsearch<br/>Log Storage)]
        JAEGER_DB[(Jaeger<br/>Trace Storage)]
    end
    
    subgraph "Visualization"
        GRAFANA[Grafana<br/>Dashboards]
        KIBANA[Kibana<br/>Log Analysis]
        JAEGER_UI[Jaeger UI<br/>Trace Analysis]
    end
    
    subgraph "Alerting"
        ALERT_MGR[AlertManager<br/>Alert Routing]
        PAGERDUTY[PagerDuty<br/>Incident Management]
        SLACK[Slack<br/>Team Notifications]
        EMAIL[Email<br/>Notifications]
    end
    
    APP --> PROM
    INFRA --> PROM
    K8S --> PROM
    LOGS --> ELASTIC
    TRACES --> JAEGER_DB
    
    PROM --> GRAFANA
    PROM --> ALERT_MGR
    ELASTIC --> KIBANA
    JAEGER_DB --> JAEGER_UI
    
    ALERT_MGR --> PAGERDUTY
    ALERT_MGR --> SLACK
    ALERT_MGR --> EMAIL
```

### 2.2 Monitoring Components

| Component | Purpose | Technology | Retention |
|-----------|---------|------------|-----------|
| **Metrics Collection** | Application & infrastructure metrics | Prometheus | 30 days |
| **Log Aggregation** | Centralized logging | ELK Stack | 90 days |
| **Distributed Tracing** | Request flow analysis | Jaeger | 7 days |
| **Visualization** | Dashboards and analysis | Grafana + Kibana | N/A |
| **Alerting** | Incident notification | AlertManager + PagerDuty | 1 year |
| **Status Page** | Public status communication | StatusPage.io | 1 year |

## 3. Service Level Indicators (SLIs) and Objectives (SLOs)

### 3.1 Primary SLIs/SLOs

| SLI | SLO Target | Measurement Window | Alert Threshold |
|-----|------------|-------------------|-----------------|
| **Availability** | ≥ 99.9% | 30 days | < 99.5% |
| **API Latency (95th percentile)** | < 200ms | 5 minutes | > 500ms |
| **API Latency (99th percentile)** | < 500ms | 5 minutes | > 1000ms |
| **Error Rate** | < 0.1% | 5 minutes | > 1% |
| **Database Query Latency (95th)** | < 100ms | 5 minutes | > 200ms |
| **Cache Hit Rate** | > 90% | 5 minutes | < 80% |

### 3.2 Business SLIs/SLOs

| SLI | SLO Target | Measurement Window | Alert Threshold |
|-----|------------|-------------------|-----------------|
| **Successful Payments** | > 99.5% | 1 hour | < 95% |
| **Tracking Link Generation** | > 99.9% | 5 minutes | < 99% |
| **User Authentication** | > 99.9% | 5 minutes | < 99% |
| **Data Sync Success** | > 99% | 1 hour | < 95% |
| **Report Generation** | > 99% | 15 minutes | < 95% |

### 3.3 SLO Monitoring Queries

#### Availability SLO
```promql
# Availability calculation (based on successful requests)
(
  sum(rate(http_requests_total{job="affiliate-api",code!~"5.."}[5m])) /
  sum(rate(http_requests_total{job="affiliate-api"}[5m]))
) * 100
```

#### Latency SLO
```promql
# 95th percentile latency
histogram_quantile(0.95, 
  sum(rate(http_request_duration_seconds_bucket{job="affiliate-api"}[5m])) by (le)
)

# 99th percentile latency
histogram_quantile(0.99, 
  sum(rate(http_request_duration_seconds_bucket{job="affiliate-api"}[5m])) by (le)
)
```

#### Error Rate SLO
```promql
# Error rate calculation
(
  sum(rate(http_requests_total{job="affiliate-api",code=~"5.."}[5m])) /
  sum(rate(http_requests_total{job="affiliate-api"}[5m]))
) * 100
```

## 4. Dashboard Configuration

### 4.1 Executive Dashboard

**Purpose**: High-level service health for management  
**Audience**: Executives, Product Owners, Service Managers  
**Update Frequency**: Real-time

```json
{
  "dashboard": {
    "title": "Affiliate Platform - Executive Overview",
    "tags": ["executive", "overview"],
    "time": {
      "from": "now-24h",
      "to": "now"
    },
    "panels": [
      {
        "title": "Service Availability (24h)",
        "type": "stat",
        "targets": [
          {
            "expr": "avg_over_time((sum(rate(http_requests_total{job=\"affiliate-api\",code!~\"5..\"}[5m])) / sum(rate(http_requests_total{job=\"affiliate-api\"}[5m])) * 100)[24h:])",
            "legendFormat": "Availability %"
          }
        ],
        "fieldConfig": {
          "defaults": {
            "unit": "percent",
            "thresholds": {
              "steps": [
                {"color": "red", "value": 0},
                {"color": "yellow", "value": 99},
                {"color": "green", "value": 99.9}
              ]
            }
          }
        }
      },
      {
        "title": "Active Users (24h)",
        "type": "stat",
        "targets": [
          {
            "expr": "count(count by (user_id) (increase(http_requests_total{job=\"affiliate-api\"}[24h])))",
            "legendFormat": "Active Users"
          }
        ]
      },
      {
        "title": "Revenue Impact",
        "type": "stat",
        "targets": [
          {
            "expr": "sum(increase(payment_processed_total{status=\"success\"}[24h]))",
            "legendFormat": "Successful Payments"
          }
        ]
      },
      {
        "title": "System Health",
        "type": "table",
        "targets": [
          {
            "expr": "up{job=\"affiliate-api\"}",
            "format": "table",
            "legendFormat": "{{instance}}"
          }
        ]
      }
    ]
  }
}
```

### 4.2 Operations Dashboard

**Purpose**: Detailed technical metrics for operations team  
**Audience**: SRE, DevOps, On-call engineers  
**Update Frequency**: Real-time

```json
{
  "dashboard": {
    "title": "Affiliate Platform - Operations",
    "tags": ["operations", "technical"],
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(http_requests_total{job=\"affiliate-api\"}[5m])) by (method, status)",
            "legendFormat": "{{method}} {{status}}"
          }
        ]
      },
      {
        "title": "Response Time Distribution",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.50, sum(rate(http_request_duration_seconds_bucket{job=\"affiliate-api\"}[5m])) by (le))",
            "legendFormat": "50th percentile"
          },
          {
            "expr": "histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket{job=\"affiliate-api\"}[5m])) by (le))",
            "legendFormat": "95th percentile"
          },
          {
            "expr": "histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket{job=\"affiliate-api\"}[5m])) by (le))",
            "legendFormat": "99th percentile"
          }
        ]
      },
      {
        "title": "Database Performance",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(postgresql_queries_total[5m])",
            "legendFormat": "Queries/sec"
          },
          {
            "expr": "postgresql_connections_active",
            "legendFormat": "Active Connections"
          }
        ]
      },
      {
        "title": "Resource Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(container_cpu_usage_seconds_total{pod=~\"affiliate-api-.*\"}[5m]) * 100",
            "legendFormat": "CPU % - {{pod}}"
          },
          {
            "expr": "container_memory_usage_bytes{pod=~\"affiliate-api-.*\"} / 1024 / 1024",
            "legendFormat": "Memory MB - {{pod}}"
          }
        ]
      }
    ]
  }
}
```

### 4.3 Business Metrics Dashboard

**Purpose**: Business KPIs and revenue metrics  
**Audience**: Product team, Business analysts  
**Update Frequency**: Hourly

```json
{
  "dashboard": {
    "title": "Affiliate Platform - Business Metrics",
    "tags": ["business", "kpi"],
    "panels": [
      {
        "title": "Daily Active Organizations",
        "type": "graph",
        "targets": [
          {
            "expr": "count(count by (organization_id) (increase(http_requests_total{job=\"affiliate-api\"}[24h])))",
            "legendFormat": "Active Organizations"
          }
        ]
      },
      {
        "title": "Campaign Performance",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(tracking_link_clicks_total[1h])) by (campaign_id)",
            "legendFormat": "Clicks/hour - Campaign {{campaign_id}}"
          }
        ]
      },
      {
        "title": "Payment Processing",
        "type": "graph",
        "targets": [
          {
            "expr": "sum(rate(payment_processed_total{status=\"success\"}[1h]))",
            "legendFormat": "Successful Payments/hour"
          },
          {
            "expr": "sum(rate(payment_processed_total{status=\"failed\"}[1h]))",
            "legendFormat": "Failed Payments/hour"
          }
        ]
      },
      {
        "title": "API Usage by Endpoint",
        "type": "table",
        "targets": [
          {
            "expr": "topk(10, sum(rate(http_requests_total{job=\"affiliate-api\"}[1h])) by (endpoint))",
            "format": "table"
          }
        ]
      }
    ]
  }
}
```

## 5. Alert Rules Configuration

### 5.1 Critical Alerts (Severity 1)

#### Service Down Alert
```yaml
groups:
- name: critical-alerts
  rules:
  - alert: ServiceDown
    expr: up{job="affiliate-api"} == 0
    for: 1m
    labels:
      severity: critical
      team: sre
      service: affiliate-api
    annotations:
      summary: "Affiliate API service is down"
      description: "The Affiliate API service has been down for more than 1 minute"
      runbook_url: "https://docs.company.com/runbooks/service-down"
      dashboard_url: "https://grafana.company.com/d/operations"
```

#### High Error Rate Alert
```yaml
  - alert: HighErrorRate
    expr: |
      (
        sum(rate(http_requests_total{job="affiliate-api",code=~"5.."}[5m])) /
        sum(rate(http_requests_total{job="affiliate-api"}[5m]))
      ) * 100 > 5
    for: 2m
    labels:
      severity: critical
      team: sre
      service: affiliate-api
    annotations:
      summary: "High error rate detected"
      description: "Error rate is {{ $value }}% for the last 5 minutes"
      runbook_url: "https://docs.company.com/runbooks/high-error-rate"
```

#### Database Connection Failure
```yaml
  - alert: DatabaseConnectionFailure
    expr: postgresql_up == 0
    for: 30s
    labels:
      severity: critical
      team: sre
      service: database
    annotations:
      summary: "Database connection failed"
      description: "Cannot connect to PostgreSQL database"
      runbook_url: "https://docs.company.com/runbooks/database-down"
```

### 5.2 Warning Alerts (Severity 2)

#### High Response Time Alert
```yaml
- name: warning-alerts
  rules:
  - alert: HighResponseTime
    expr: |
      histogram_quantile(0.95, 
        sum(rate(http_request_duration_seconds_bucket{job="affiliate-api"}[5m])) by (le)
      ) > 0.5
    for: 5m
    labels:
      severity: warning
      team: sre
      service: affiliate-api
    annotations:
      summary: "High response time detected"
      description: "95th percentile response time is {{ $value }}s"
      runbook_url: "https://docs.company.com/runbooks/high-latency"
```

#### High CPU Usage Alert
```yaml
  - alert: HighCPUUsage
    expr: |
      (
        rate(container_cpu_usage_seconds_total{pod=~"affiliate-api-.*"}[5m]) * 100
      ) > 80
    for: 10m
    labels:
      severity: warning
      team: sre
      service: affiliate-api
    annotations:
      summary: "High CPU usage on {{ $labels.pod }}"
      description: "CPU usage is {{ $value }}% on pod {{ $labels.pod }}"
      runbook_url: "https://docs.company.com/runbooks/high-cpu"
```

#### Low Cache Hit Rate Alert
```yaml
  - alert: LowCacheHitRate
    expr: |
      (
        redis_keyspace_hits_total / 
        (redis_keyspace_hits_total + redis_keyspace_misses_total)
      ) * 100 < 80
    for: 15m
    labels:
      severity: warning
      team: sre
      service: redis
    annotations:
      summary: "Low Redis cache hit rate"
      description: "Cache hit rate is {{ $value }}%"
      runbook_url: "https://docs.company.com/runbooks/low-cache-hit-rate"
```

### 5.3 Business Alerts

#### Payment Processing Failure
```yaml
- name: business-alerts
  rules:
  - alert: PaymentProcessingFailure
    expr: |
      (
        sum(rate(payment_processed_total{status="failed"}[5m])) /
        sum(rate(payment_processed_total[5m]))
      ) * 100 > 10
    for: 5m
    labels:
      severity: critical
      team: business
      service: payments
    annotations:
      summary: "High payment failure rate"
      description: "Payment failure rate is {{ $value }}%"
      runbook_url: "https://docs.company.com/runbooks/payment-failures"
```

#### Low User Activity Alert
```yaml
  - alert: LowUserActivity
    expr: |
      sum(rate(http_requests_total{job="affiliate-api"}[1h])) < 10
    for: 30m
    labels:
      severity: warning
      team: business
      service: affiliate-api
    annotations:
      summary: "Unusually low user activity"
      description: "Request rate is {{ $value }} requests/second"
      runbook_url: "https://docs.company.com/runbooks/low-activity"
```

## 6. Alert Routing and Escalation

### 6.1 AlertManager Configuration

```yaml
global:
  smtp_smarthost: 'smtp.company.com:587'
  smtp_from: 'alerts@company.com'
  pagerduty_url: 'https://events.pagerduty.com/v2/enqueue'

route:
  group_by: ['alertname', 'service']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 12h
  receiver: 'default'
  routes:
  - match:
      severity: critical
    receiver: 'pagerduty-critical'
    routes:
    - match:
        team: sre
      receiver: 'pagerduty-sre'
    - match:
        team: business
      receiver: 'pagerduty-business'
  - match:
      severity: warning
    receiver: 'slack-warnings'

receivers:
- name: 'default'
  slack_configs:
  - api_url: 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX'
    channel: '#alerts'
    title: 'Alert: {{ .GroupLabels.alertname }}'
    text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'

- name: 'pagerduty-critical'
  pagerduty_configs:
  - routing_key: 'PAGERDUTY_INTEGRATION_KEY'
    description: '{{ .GroupLabels.alertname }}: {{ .GroupLabels.service }}'
    details:
      firing: '{{ .Alerts.Firing | len }}'
      resolved: '{{ .Alerts.Resolved | len }}'

- name: 'pagerduty-sre'
  pagerduty_configs:
  - routing_key: 'PAGERDUTY_SRE_KEY'
    description: 'SRE Alert: {{ .GroupLabels.alertname }}'

- name: 'slack-warnings'
  slack_configs:
  - api_url: 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX'
    channel: '#sre-warnings'
    title: 'Warning: {{ .GroupLabels.alertname }}'
    text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
```

### 6.2 Escalation Matrix

| Alert Severity | Initial Notification | Escalation (15 min) | Escalation (30 min) | Escalation (1 hour) |
|----------------|---------------------|-------------------|-------------------|-------------------|
| **Critical** | On-call SRE (PagerDuty) | SRE Lead | Service Manager | CTO |
| **Warning** | Slack #sre-warnings | On-call SRE | SRE Lead | - |
| **Business Critical** | Business On-call | Product Owner | VP Product | CEO |
| **Security** | Security Team | CISO | CTO | CEO |

### 6.3 On-Call Schedule

#### Primary On-Call Rotation
```yaml
SRE Team Rotation (Weekly):
  Week 1: Alex Rodriguez (SRE Lead)
  Week 2: Senior SRE Engineer 1
  Week 3: Senior SRE Engineer 2
  Week 4: SRE Engineer 1

Business Team Rotation (Weekly):
  Week 1: Product Owner
  Week 2: Customer Success Manager
  Week 3: Business Analyst
  Week 4: Product Manager

Security Team Rotation (Monthly):
  Month 1: Security Officer
  Month 2: Security Engineer
  Month 3: CISO
```

#### Escalation Contacts
```yaml
Primary Escalation:
  - SRE Lead: Alex Rodriguez (+1-555-0103)
  - Service Manager: Mike Chen (+1-555-0102)
  - Product Owner: Sarah Johnson (+1-555-0101)

Secondary Escalation:
  - VP Engineering: (+1-555-0301)
  - VP Product: (+1-555-0302)
  - CTO: (+1-555-0303)

Emergency Escalation:
  - CEO: (+1-555-0401)
  - Board Chair: (+1-555-0402)
```

## 7. Log Management and Analysis

### 7.1 Log Structure and Standards

#### Application Log Format
```json
{
  "timestamp": "2025-08-05T14:30:00.123Z",
  "level": "INFO",
  "service": "affiliate-api",
  "version": "v1.2.3",
  "request_id": "req-123456789",
  "user_id": "user-987654321",
  "organization_id": "org-456789123",
  "method": "POST",
  "endpoint": "/api/v1/campaigns",
  "status_code": 201,
  "response_time_ms": 145,
  "message": "Campaign created successfully",
  "metadata": {
    "campaign_id": "camp-123456",
    "advertiser_id": "adv-789123"
  }
}
```

#### Error Log Format
```json
{
  "timestamp": "2025-08-05T14:30:00.123Z",
  "level": "ERROR",
  "service": "affiliate-api",
  "version": "v1.2.3",
  "request_id": "req-123456789",
  "error": {
    "type": "DatabaseConnectionError",
    "message": "Connection to database failed",
    "stack_trace": "...",
    "code": "DB_CONN_001"
  },
  "context": {
    "user_id": "user-987654321",
    "endpoint": "/api/v1/campaigns",
    "retry_count": 3
  }
}
```

### 7.2 Log Analysis Queries

#### Common Error Analysis
```json
{
  "query": {
    "bool": {
      "must": [
        {"term": {"level": "ERROR"}},
        {"range": {"timestamp": {"gte": "now-1h"}}}
      ]
    }
  },
  "aggs": {
    "error_types": {
      "terms": {
        "field": "error.type.keyword",
        "size": 10
      }
    },
    "error_endpoints": {
      "terms": {
        "field": "endpoint.keyword",
        "size": 10
      }
    }
  }
}
```

#### Performance Analysis
```json
{
  "query": {
    "bool": {
      "must": [
        {"range": {"response_time_ms": {"gte": 1000}}},
        {"range": {"timestamp": {"gte": "now-1h"}}}
      ]
    }
  },
  "aggs": {
    "slow_endpoints": {
      "terms": {
        "field": "endpoint.keyword",
        "size": 10
      },
      "aggs": {
        "avg_response_time": {
          "avg": {"field": "response_time_ms"}
        }
      }
    }
  }
}
```

#### User Activity Analysis
```json
{
  "query": {
    "bool": {
      "must": [
        {"exists": {"field": "user_id"}},
        {"range": {"timestamp": {"gte": "now-24h"}}}
      ]
    }
  },
  "aggs": {
    "active_users": {
      "cardinality": {"field": "user_id.keyword"}
    },
    "user_activity": {
      "terms": {
        "field": "user_id.keyword",
        "size": 100
      },
      "aggs": {
        "request_count": {
          "value_count": {"field": "request_id.keyword"}
        }
      }
    }
  }
}
```

### 7.3 Log Retention and Archival

| Log Type | Retention Period | Storage Tier | Compression |
|----------|------------------|--------------|-------------|
| **Application Logs** | 90 days | Hot (0-7 days), Warm (8-30 days), Cold (31-90 days) | gzip |
| **Error Logs** | 1 year | Hot (0-30 days), Cold (31-365 days) | gzip |
| **Audit Logs** | 7 years | Hot (0-90 days), Cold (91 days-7 years) | gzip |
| **Access Logs** | 1 year | Warm (0-30 days), Cold (31-365 days) | gzip |
| **Debug Logs** | 7 days | Hot only | gzip |

## 8. Distributed Tracing

### 8.1 Trace Configuration

#### Jaeger Configuration
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: jaeger-config
data:
  jaeger.yaml: |
    sampling:
      default_strategy:
        type: probabilistic
        param: 0.1  # Sample 10% of traces
    storage:
      type: elasticsearch
      options:
        es:
          server-urls: http://elasticsearch:9200
          index-prefix: jaeger
    query:
      base-path: /jaeger
```

#### Application Tracing
```go
// Example Go tracing implementation
func (h *Handler) CreateCampaign(c *gin.Context) {
    // Start trace span
    span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "create_campaign")
    defer span.Finish()
    
    // Add tags
    span.SetTag("user_id", c.GetString("user_id"))
    span.SetTag("organization_id", c.GetString("organization_id"))
    
    // Database operation with tracing
    dbSpan, dbCtx := opentracing.StartSpanFromContext(ctx, "database_insert")
    campaign, err := h.campaignService.Create(dbCtx, request)
    dbSpan.Finish()
    
    if err != nil {
        span.SetTag("error", true)
        span.LogFields(log.Error(err))
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }
    
    // External API call with tracing
    extSpan, extCtx := opentracing.StartSpanFromContext(ctx, "everflow_sync")
    err = h.everflowService.SyncCampaign(extCtx, campaign)
    extSpan.Finish()
    
    span.SetTag("campaign_id", campaign.ID)
    c.JSON(201, campaign)
}
```

### 8.2 Trace Analysis

#### Common Trace Queries
```json
{
  "query": {
    "bool": {
      "must": [
        {"term": {"process.serviceName": "affiliate-api"}},
        {"range": {"startTime": {"gte": "now-1h"}}},
        {"range": {"duration": {"gte": 1000000}}}
      ]
    }
  },
  "sort": [
    {"duration": {"order": "desc"}}
  ],
  "size": 100
}
```

#### Performance Bottleneck Analysis
```json
{
  "aggs": {
    "operations": {
      "terms": {
        "field": "operationName.keyword",
        "size": 20
      },
      "aggs": {
        "avg_duration": {
          "avg": {"field": "duration"}
        },
        "max_duration": {
          "max": {"field": "duration"}
        },
        "error_rate": {
          "filter": {
            "term": {"tags.error": true}
          }
        }
      }
    }
  }
}
```

## 9. Synthetic Monitoring

### 9.1 Health Check Monitoring

#### External Health Checks
```bash
#!/bin/bash
# External synthetic monitoring script

ENDPOINTS=(
    "https://api.affiliate-platform.com/health"
    "https://api.affiliate-platform.com/api/v1/public/status"
)

for endpoint in "${ENDPOINTS[@]}"; do
    echo "Checking $endpoint"
    
    # Measure response time and status
    response=$(curl -w "%{http_code},%{time_total}" -o /dev/null -s "$endpoint")
    http_code=$(echo $response | cut -d',' -f1)
    response_time=$(echo $response | cut -d',' -f2)
    
    # Send metrics to monitoring system
    echo "http_check_status{endpoint=\"$endpoint\"} $http_code" | \
        curl -X POST --data-binary @- \
        http://pushgateway:9091/metrics/job/synthetic_monitoring
    
    echo "http_check_duration{endpoint=\"$endpoint\"} $response_time" | \
        curl -X POST --data-binary @- \
        http://pushgateway:9091/metrics/job/synthetic_monitoring
done
```

#### API Functionality Tests
```bash
#!/bin/bash
# API functionality synthetic tests

API_BASE="https://api.affiliate-platform.com"
TEST_TOKEN="$SYNTHETIC_TEST_TOKEN"

# Test authentication
auth_response=$(curl -s -w "%{http_code}" -H "Authorization: Bearer $TEST_TOKEN" \
    "$API_BASE/api/v1/users/me")
auth_status=$(echo $auth_response | tail -c 4)

if [ "$auth_status" = "200" ]; then
    echo "auth_test_success 1"
else
    echo "auth_test_success 0"
fi

# Test campaign creation
campaign_data='{"name":"Synthetic Test Campaign","status":"draft"}'
campaign_response=$(curl -s -w "%{http_code}" -X POST \
    -H "Authorization: Bearer $TEST_TOKEN" \
    -H "Content-Type: application/json" \
    -d "$campaign_data" \
    "$API_BASE/api/v1/campaigns")
campaign_status=$(echo $campaign_response | tail -c 4)

if [ "$campaign_status" = "201" ]; then
    echo "campaign_create_test_success 1"
else
    echo "campaign_create_test_success 0"
fi

# Send metrics to Prometheus
echo "synthetic_test_auth_success $auth_success" | \
    curl -X POST --data-binary @- \
    http://pushgateway:9091/metrics/job/synthetic_tests

echo "synthetic_test_campaign_create_success $campaign_success" | \
    curl -X POST --data-binary @- \
    http://pushgateway:9091/metrics/job/synthetic_tests
```

### 9.2 User Journey Monitoring

#### Critical User Journey Tests
```javascript
// Playwright synthetic test for critical user journey
const { test, expect } = require('@playwright/test');

test('Complete affiliate signup and campaign creation', async ({ page }) => {
  // Start timing
  const startTime = Date.now();
  
  // Navigate to signup page
  await page.goto('https://app.affiliate-platform.com/signup');
  
  // Fill signup form
  await page.fill('[data-testid="email"]', 'synthetic-test@example.com');
  await page.fill('[data-testid="password"]', 'SyntheticTest123!');
  await page.fill('[data-testid="organization-name"]', 'Synthetic Test Org');
  
  // Submit signup
  await page.click('[data-testid="signup-button"]');
  
  // Wait for dashboard
  await expect(page.locator('[data-testid="dashboard"]')).toBeVisible();
  
  // Create campaign
  await page.click('[data-testid="create-campaign"]');
  await page.fill('[data-testid="campaign-name"]', 'Synthetic Test Campaign');
  await page.selectOption('[data-testid="campaign-type"]', 'cpc');
  await page.fill('[data-testid="payout-amount"]', '5.00');
  
  // Submit campaign
  await page.click('[data-testid="save-campaign"]');
  
  // Verify campaign created
  await expect(page.locator('[data-testid="campaign-success"]')).toBeVisible();
  
  // Calculate total time
  const totalTime = Date.now() - startTime;
  
  // Send metrics
  await fetch('http://pushgateway:9091/metrics/job/user_journey_tests', {
    method: 'POST',
    body: `user_journey_signup_duration_ms ${totalTime}\nuser_journey_signup_success 1`
  });
});
```

## 10. Capacity Planning and Forecasting

### 10.1 Resource Usage Trends

#### CPU Usage Forecasting
```promql
# Predict CPU usage for next 7 days based on linear regression
predict_linear(
  avg_over_time(
    rate(container_cpu_usage_seconds_total{pod=~"affiliate-api-.*"}[5m])[7d:]
  ), 
  7*24*3600
) * 100
```

#### Memory Usage Forecasting
```promql
# Predict memory usage growth
predict_linear(
  avg_over_time(
    container_memory_usage_bytes{pod=~"affiliate-api-.*"}[5m])[7d:]
  ), 
  7*24*3600
) / 1024 / 1024 / 1024
```

#### Database Growth Forecasting
```promql
# Predict database size growth
predict_linear(
  postgresql_database_size_bytes[30d], 
  30*24*3600
) / 1024 / 1024 / 1024
```

### 10.2 Capacity Alerts

#### Resource Exhaustion Alerts
```yaml
- name: capacity-alerts
  rules:
  - alert: CPUCapacityWarning
    expr: |
      predict_linear(
        avg_over_time(rate(container_cpu_usage_seconds_total{pod=~"affiliate-api-.*"}[5m])[7d:]), 
        7*24*3600
      ) * 100 > 80
    for: 1h
    labels:
      severity: warning
      team: sre
    annotations:
      summary: "CPU capacity will be exhausted in 7 days"
      description: "Predicted CPU usage in 7 days: {{ $value }}%"

  - alert: MemoryCapacityWarning
    expr: |
      predict_linear(
        avg_over_time(container_memory_usage_bytes{pod=~"affiliate-api-.*"}[5m])[7d:]), 
        7*24*3600
      ) / container_spec_memory_limit_bytes > 0.9
    for: 1h
    labels:
      severity: warning
      team: sre
    annotations:
      summary: "Memory capacity will be exhausted in 7 days"
      description: "Predicted memory usage in 7 days: {{ $value }}%"

  - alert: DatabaseCapacityWarning
    expr: |
      predict_linear(postgresql_database_size_bytes[30d], 30*24*3600) / 
      postgresql_database_size_limit_bytes > 0.8
    for: 6h
    labels:
      severity: warning
      team: sre
    annotations:
      summary: "Database will reach 80% capacity in 30 days"
      description: "Predicted database usage: {{ $value }}%"
```

### 10.3 Scaling Recommendations

#### Auto-Scaling Configuration
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: affiliate-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: affiliate-api
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  - type: Pods
    pods:
      metric:
        name: http_requests_per_second
      target:
        type: AverageValue
        averageValue: "100"
  behavior:
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 60
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
```

## 11. Incident Response Integration

### 11.1 Automated Incident Creation

#### PagerDuty Integration
```yaml
# AlertManager configuration for PagerDuty
- name: 'pagerduty-integration'
  pagerduty_configs:
  - routing_key: 'PAGERDUTY_INTEGRATION_KEY'
    description: '{{ .GroupLabels.alertname }}: {{ .CommonAnnotations.summary }}'
    details:
      alert_count: '{{ .Alerts | len }}'
      service: '{{ .GroupLabels.service }}'
      severity: '{{ .GroupLabels.severity }}'
      runbook: '{{ .CommonAnnotations.runbook_url }}'
      dashboard: '{{ .CommonAnnotations.dashboard_url }}'
    links:
    - href: '{{ .CommonAnnotations.runbook_url }}'
      text: 'Runbook'
    - href: '{{ .CommonAnnotations.dashboard_url }}'
      text: 'Dashboard'
```

#### Slack Integration
```yaml
- name: 'slack-integration'
  slack_configs:
  - api_url: 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX'
    channel: '#incidents'
    title: '{{ .GroupLabels.severity | toUpper }}: {{ .GroupLabels.alertname }}'
    text: |
      {{ range .Alerts }}
      *Alert:* {{ .Annotations.summary }}
      *Description:* {{ .Annotations.description }}
      *Service:* {{ .Labels.service }}
      *Runbook:* {{ .Annotations.runbook_url }}
      {{ end }}
    actions:
    - type: button
      text: 'View Dashboard'
      url: '{{ .CommonAnnotations.dashboard_url }}'
    - type: button
      text: 'View Runbook'
      url: '{{ .CommonAnnotations.runbook_url }}'
```

### 11.2 Status Page Integration

#### Automated Status Updates
```bash
#!/bin/bash
# Automated status page update based on alerts

ALERT_NAME="$1"
ALERT_STATUS="$2"  # firing or resolved
COMPONENT_ID="$3"

case "$ALERT_NAME" in
    "ServiceDown")
        STATUS="major_outage"
        MESSAGE="We are investigating reports of service unavailability."
        ;;
    "HighErrorRate")
        STATUS="partial_outage"
        MESSAGE="Some users may experience errors. We are investigating."
        ;;
    "HighResponseTime")
        STATUS="degraded_performance"
        MESSAGE="Users may experience slower than normal response times."
        ;;
    *)
        STATUS="operational"
        MESSAGE="All systems operational."
        ;;
esac

if [ "$ALERT_STATUS" = "resolved" ]; then
    STATUS="operational"
    MESSAGE="Service has been restored to normal operation."
fi

# Update component status
curl -X PATCH "https://api.statuspage.io/v1/pages/$STATUSPAGE_PAGE_ID/components/$COMPONENT_ID" \
    -H "Authorization: OAuth $STATUSPAGE_API_KEY" \
    -H "Content-Type: application/json" \
    -d "{\"component\": {\"status\": \"$STATUS\"}}"

# Create or update incident
if [ "$STATUS" != "operational" ] && [ "$ALERT_STATUS" = "firing" ]; then
    curl -X POST "https://api.statuspage.io/v1/pages/$STATUSPAGE_PAGE_ID/incidents" \
        -H "Authorization: OAuth $STATUSPAGE_API_KEY" \
        -H "Content-Type: application/json" \
        -d "{
            \"incident\": {
                \"name\": \"$ALERT_NAME\",
                \"status\": \"investigating\",
                \"impact_override\": \"major\",
                \"body\": \"$MESSAGE\"
            }
        }"
fi
```

---

## Appendix A: Monitoring Tools Quick Reference

### Prometheus Queries
```promql
# Service availability
up{job="affiliate-api"}

# Request rate
rate(http_requests_total{job="affiliate-api"}[5m])

# Error rate
rate(http_requests_total{job="affiliate-api",code=~"5.."}[5m]) / 
rate(http_requests_total{job="affiliate-api"}[5m])

# Response time percentiles
histogram_quantile(0.95, 
  rate(http_request_duration_seconds_bucket{job="affiliate-api"}[5m])
)

# Database connections
postgresql_connections_active

# Memory usage
container_memory_usage_bytes{pod=~"affiliate-api-.*"}

# CPU usage
rate(container_cpu_usage_seconds_total{pod=~"affiliate-api-.*"}[5m])
```

### Grafana Dashboard URLs
- **Executive Dashboard**: https://grafana.company.com/d/executive
- **Operations Dashboard**: https://grafana.company.com/d/operations  
- **Business Metrics**: https://grafana.company.com/d/business
- **Infrastructure**: https://grafana.company.com/d/infrastructure
- **Database**: https://grafana.company.com/d/database

### Kibana Saved Searches
- **Error Analysis**: https://kibana.company.com/app/discover#/error-analysis
- **Performance Issues**: https://kibana.company.com/app/discover#/performance
- **User Activity**: https://kibana.company.com/app/discover#/user-activity
- **Security Events**: https://kibana.company.com/app/discover#/security

## Appendix B: Alert Runbook Links

| Alert Name | Runbook URL | Dashboard URL |
|------------|-------------|---------------|
| ServiceDown | https://docs.company.com/runbooks/service-down | https://grafana.company.com/d/operations |
| HighErrorRate | https://docs.company.com/runbooks/high-error-rate | https://grafana.company.com/d/operations |
| HighResponseTime | https://docs.company.com/runbooks/high-latency | https://grafana.company.com/d/operations |
| DatabaseConnectionFailure | https://docs.company.com/runbooks/database-down | https://grafana.company.com/d/database |
| HighCPUUsage | https://docs.company.com/runbooks/high-cpu | https://grafana.company.com/d/infrastructure |
| PaymentProcessingFailure | https://docs.company.com/runbooks/payment-failures | https://grafana.company.com/d/business |

---

**Document Classification**: Operational  
**Access Level**: SRE Team, Operations Team, On-Call Engineers  
**Review Frequency**: Monthly  
**Related Documents**: Runbook, System Architecture, SLA