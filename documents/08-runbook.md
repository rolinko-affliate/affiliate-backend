# Runbook: Affiliate Backend Platform

**Document Version**: v1.0  
**Owner**: On-Call Team Lead  
**Last Updated**: 2025-08-05  
**Next Review**: 2026-02-05

---

## 1. Overview

This runbook provides step-by-step operational procedures for the Affiliate Backend Platform. It is designed for on-call engineers and operations teams to quickly diagnose and resolve common issues, perform routine maintenance, and execute emergency procedures.

### Quick Reference
- **Service URL**: https://api.affiliate-platform.com
- **Status Page**: https://status.affiliate-platform.com
- **Monitoring**: https://monitoring.affiliate-platform.com
- **Emergency Escalation**: +1-555-0199

### Service Architecture Quick View
```
Internet → Cloudflare → GCP Load Balancer → GKE Cluster → PostgreSQL
                                        ↓
                                   Redis Cache
```

## 2. Daily Operations Checklist

### 2.1 Morning Health Check (9:00 AM EST)

```bash
#!/bin/bash
# Daily morning health check script
# Run this every morning to verify system health

echo "=== Daily Health Check - $(date) ==="

# 1. Check service availability
echo "1. Checking service availability..."
if curl -f -s https://api.affiliate-platform.com/health > /dev/null; then
    echo "✅ API is responding"
else
    echo "❌ API is not responding - INVESTIGATE IMMEDIATELY"
fi

# 2. Check database connectivity
echo "2. Checking database connectivity..."
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    echo 'SELECT 1;' | psql \$DATABASE_URL -t
" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Database is accessible"
else
    echo "❌ Database connection failed - CHECK DATABASE STATUS"
fi

# 3. Check Redis connectivity
echo "3. Checking Redis connectivity..."
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    redis-cli -h \$REDIS_HOST ping
" | grep -q "PONG"
if [ $? -eq 0 ]; then
    echo "✅ Redis is responding"
else
    echo "❌ Redis connection failed - CHECK REDIS STATUS"
fi

# 4. Check pod status
echo "4. Checking pod status..."
UNHEALTHY_PODS=$(kubectl get pods -l app=affiliate-api --no-headers | grep -v Running | wc -l)
if [ $UNHEALTHY_PODS -eq 0 ]; then
    echo "✅ All pods are healthy"
else
    echo "❌ $UNHEALTHY_PODS unhealthy pods found - CHECK POD STATUS"
    kubectl get pods -l app=affiliate-api
fi

# 5. Check recent errors
echo "5. Checking recent errors (last 1 hour)..."
ERROR_COUNT=$(kubectl logs -l app=affiliate-api --since=1h | grep -i error | wc -l)
if [ $ERROR_COUNT -lt 10 ]; then
    echo "✅ Error count is normal ($ERROR_COUNT errors)"
else
    echo "⚠️  High error count: $ERROR_COUNT errors in last hour - REVIEW LOGS"
fi

# 6. Check disk usage
echo "6. Checking disk usage..."
kubectl top nodes | awk 'NR>1 {if ($5+0 > 80) print "❌ High disk usage on " $1 ": " $5; else print "✅ Disk usage OK on " $1 ": " $5}'

# 7. Check certificate expiry
echo "7. Checking certificate expiry..."
CERT_DAYS=$(echo | openssl s_client -servername api.affiliate-platform.com -connect api.affiliate-platform.com:443 2>/dev/null | openssl x509 -noout -dates | grep notAfter | cut -d= -f2 | xargs -I {} date -d {} +%s)
CURRENT_DATE=$(date +%s)
DAYS_UNTIL_EXPIRY=$(( (CERT_DAYS - CURRENT_DATE) / 86400 ))

if [ $DAYS_UNTIL_EXPIRY -gt 30 ]; then
    echo "✅ Certificate expires in $DAYS_UNTIL_EXPIRY days"
else
    echo "⚠️  Certificate expires in $DAYS_UNTIL_EXPIRY days - RENEW SOON"
fi

echo "=== Health Check Complete ==="
```

### 2.2 Weekly Maintenance Tasks (Sundays 2:00 AM EST)

```bash
#!/bin/bash
# Weekly maintenance script
# Run every Sunday during low-traffic period

echo "=== Weekly Maintenance - $(date) ==="

# 1. Database maintenance
echo "1. Running database maintenance..."
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c 'VACUUM ANALYZE;'
    psql \$DATABASE_URL -c 'REINDEX DATABASE affiliate_platform;'
"

# 2. Clear old logs
echo "2. Clearing old application logs..."
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    find /var/log -name '*.log' -mtime +7 -delete
"

# 3. Update container images (if auto-update enabled)
echo "3. Checking for container updates..."
kubectl rollout status deployment/affiliate-api

# 4. Backup verification
echo "4. Verifying recent backups..."
LATEST_BACKUP=$(gcloud sql backups list --instance=affiliate-prod-db --limit=1 --format="value(startTime)")
BACKUP_AGE=$(( ($(date +%s) - $(date -d "$LATEST_BACKUP" +%s)) / 3600 ))

if [ $BACKUP_AGE -lt 25 ]; then
    echo "✅ Latest backup is $BACKUP_AGE hours old"
else
    echo "❌ Latest backup is $BACKUP_AGE hours old - CHECK BACKUP SYSTEM"
fi

# 5. Security scan
echo "5. Running security scan..."
trivy image gcr.io/affiliate-platform-prod/affiliate-api:latest --severity HIGH,CRITICAL --quiet

# 6. Performance metrics review
echo "6. Generating performance report..."
kubectl top pods -l app=affiliate-api

echo "=== Weekly Maintenance Complete ==="
```

## 3. Common Issue Troubleshooting

### 3.1 Service Unavailable (HTTP 503)

**Symptoms**: API returns 503 errors, health check fails

**Diagnosis Steps**:
```bash
# 1. Check pod status
kubectl get pods -l app=affiliate-api

# 2. Check pod logs for errors
kubectl logs -l app=affiliate-api --tail=50

# 3. Check resource usage
kubectl top pods -l app=affiliate-api

# 4. Check node status
kubectl get nodes
kubectl describe nodes
```

**Common Causes & Solutions**:

#### Cause 1: Pods are not ready
```bash
# Check pod readiness
kubectl describe pods -l app=affiliate-api

# If pods are failing readiness checks:
# 1. Check if database is accessible
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "psql \$DATABASE_URL -c 'SELECT 1;'"

# 2. Restart pods if needed
kubectl rollout restart deployment/affiliate-api
```

#### Cause 2: Resource exhaustion
```bash
# Check resource limits
kubectl describe deployment affiliate-api

# If CPU/Memory limits are hit:
# 1. Scale up replicas temporarily
kubectl scale deployment affiliate-api --replicas=5

# 2. Increase resource limits (requires deployment update)
kubectl patch deployment affiliate-api -p '{"spec":{"template":{"spec":{"containers":[{"name":"api","resources":{"limits":{"memory":"1Gi","cpu":"1000m"}}}]}}}}'
```

#### Cause 3: Database connectivity issues
```bash
# Check database status
gcloud sql instances describe affiliate-prod-db

# Check database connections
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c 'SELECT count(*) FROM pg_stat_activity;'
"

# If connection pool exhausted, restart pods
kubectl rollout restart deployment/affiliate-api
```

### 3.2 High Response Times (> 1 second)

**Symptoms**: API responses are slow, timeout errors

**Diagnosis Steps**:
```bash
# 1. Check current response times
curl -w "@curl-format.txt" -o /dev/null -s https://api.affiliate-platform.com/health

# Create curl-format.txt:
cat > curl-format.txt << 'EOF'
     time_namelookup:  %{time_namelookup}\n
        time_connect:  %{time_connect}\n
     time_appconnect:  %{time_appconnect}\n
    time_pretransfer:  %{time_pretransfer}\n
       time_redirect:  %{time_redirect}\n
  time_starttransfer:  %{time_starttransfer}\n
                     ----------\n
          time_total:  %{time_total}\n
EOF

# 2. Check database query performance
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT query, mean_time, calls 
        FROM pg_stat_statements 
        ORDER BY mean_time DESC 
        LIMIT 10;
    \"
"

# 3. Check Redis performance
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    redis-cli -h \$REDIS_HOST info stats | grep instantaneous
"
```

**Common Solutions**:

#### Solution 1: Database optimization
```bash
# Check for long-running queries
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT pid, now() - pg_stat_activity.query_start AS duration, query 
        FROM pg_stat_activity 
        WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';
    \"
"

# Kill long-running queries if necessary
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c 'SELECT pg_terminate_backend(PID);'
"

# Update table statistics
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c 'ANALYZE;'
"
```

#### Solution 2: Cache warming
```bash
# Warm up Redis cache
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    curl -X POST http://localhost:8080/internal/cache/warm
"

# Check cache hit rate
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    redis-cli -h \$REDIS_HOST info stats | grep keyspace
"
```

#### Solution 3: Scale resources
```bash
# Scale up pods
kubectl scale deployment affiliate-api --replicas=6

# Check if scaling helped
kubectl top pods -l app=affiliate-api
```

### 3.3 Database Connection Errors

**Symptoms**: "connection refused", "too many connections"

**Diagnosis Steps**:
```bash
# 1. Check database status
gcloud sql instances describe affiliate-prod-db

# 2. Check connection count
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT count(*) as active_connections, 
               max_conn, 
               max_conn - count(*) as available_connections
        FROM pg_stat_activity, 
             (SELECT setting::int as max_conn FROM pg_settings WHERE name='max_connections') mc;
    \"
"

# 3. Check for connection leaks
kubectl logs -l app=affiliate-api | grep -i "connection"
```

**Solutions**:

#### Solution 1: Connection pool tuning
```bash
# Check current pool settings
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    echo 'Current pool size:' \$DB_POOL_SIZE
    echo 'Max connections:' \$DB_MAX_CONNECTIONS
"

# Restart pods to reset connection pools
kubectl rollout restart deployment/affiliate-api
```

#### Solution 2: Kill idle connections
```bash
# Kill idle connections older than 1 hour
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT pg_terminate_backend(pid)
        FROM pg_stat_activity
        WHERE state = 'idle'
        AND state_change < now() - interval '1 hour';
    \"
"
```

### 3.4 Redis Connection Issues

**Symptoms**: Cache misses, Redis timeout errors

**Diagnosis Steps**:
```bash
# 1. Check Redis status
gcloud redis instances describe affiliate-prod-redis --region=us-central1

# 2. Test Redis connectivity
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    redis-cli -h \$REDIS_HOST ping
"

# 3. Check Redis memory usage
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    redis-cli -h \$REDIS_HOST info memory
"
```

**Solutions**:

#### Solution 1: Clear Redis cache
```bash
# Clear all cache (use with caution)
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    redis-cli -h \$REDIS_HOST flushall
"

# Clear specific keys
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    redis-cli -h \$REDIS_HOST --scan --pattern 'cache:*' | xargs redis-cli -h \$REDIS_HOST del
"
```

#### Solution 2: Restart Redis connection
```bash
# Restart application pods to reset Redis connections
kubectl rollout restart deployment/affiliate-api
```

## 4. Incident Response Playbooks

### 4.1 Severity 1: Complete Service Outage

**Definition**: API completely unavailable, all users affected

**Immediate Actions (0-15 minutes)**:
```bash
# 1. Acknowledge the incident
echo "Severity 1 incident acknowledged at $(date)"

# 2. Check overall system status
kubectl get pods,services,ingress -A

# 3. Check external dependencies
curl -I https://api.supabase.com/health
curl -I https://api.stripe.com/v1

# 4. Notify incident response team
# Send to Slack: #incident-response
# Page on-call: Use PagerDuty escalation

# 5. Update status page
curl -X POST "https://api.statuspage.io/v1/pages/$PAGE_ID/incidents" \
    -H "Authorization: OAuth $STATUSPAGE_API_KEY" \
    -d '{
        "incident": {
            "name": "Service Outage",
            "status": "investigating",
            "impact_override": "major",
            "body": "We are investigating reports of service unavailability."
        }
    }'
```

**Investigation Steps (15-30 minutes)**:
```bash
# 1. Check infrastructure health
gcloud compute instances list
gcloud sql instances list
kubectl get nodes

# 2. Check recent deployments
kubectl rollout history deployment/affiliate-api

# 3. Check logs for errors
kubectl logs -l app=affiliate-api --since=1h | grep -i error | tail -20

# 4. Check monitoring dashboards
# - CPU/Memory usage
# - Database performance
# - Network connectivity
```

**Resolution Steps**:
```bash
# Option 1: Rollback recent deployment
kubectl rollout undo deployment/affiliate-api

# Option 2: Scale up resources
kubectl scale deployment affiliate-api --replicas=6

# Option 3: Restart services
kubectl rollout restart deployment/affiliate-api

# Option 4: Failover to DR (if regional issue)
# Follow DR procedures in BCP/DR Plan
```

### 4.2 Severity 2: Performance Degradation

**Definition**: Service is slow but functional, some users affected

**Actions**:
```bash
# 1. Identify performance bottleneck
kubectl top pods -l app=affiliate-api
kubectl top nodes

# 2. Check database performance
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT query, mean_time, calls 
        FROM pg_stat_statements 
        WHERE mean_time > 1000 
        ORDER BY mean_time DESC 
        LIMIT 5;
    \"
"

# 3. Scale resources if needed
kubectl scale deployment affiliate-api --replicas=5

# 4. Monitor improvement
watch kubectl top pods -l app=affiliate-api
```

### 4.3 Severity 3: Individual User Issues

**Definition**: Specific user or feature affected, limited impact

**Actions**:
```bash
# 1. Identify affected user/feature
# Check logs for specific user ID or feature
kubectl logs -l app=affiliate-api | grep "user_id:$USER_ID"

# 2. Check user-specific data
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT * FROM profiles WHERE id = '$USER_ID';
    \"
"

# 3. Clear user-specific cache
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    redis-cli -h \$REDIS_HOST del user:$USER_ID:*
"

# 4. Test user functionality
curl -H "Authorization: Bearer $USER_TOKEN" \
     https://api.affiliate-platform.com/api/v1/users/me
```

## 5. Deployment Procedures

### 5.1 Standard Deployment

**Pre-deployment Checklist**:
```bash
# 1. Verify staging deployment
kubectl config use-context affiliate-staging-gke
kubectl get pods -l app=affiliate-api

# 2. Run smoke tests
curl -f https://staging-api.affiliate-platform.com/health

# 3. Check for breaking changes
git log --oneline HEAD~5..HEAD

# 4. Verify database migrations (if any)
kubectl exec -it deployment/affiliate-api -- ./migrate status
```

**Deployment Steps**:
```bash
# 1. Switch to production context
kubectl config use-context affiliate-prod-gke

# 2. Update deployment with new image
kubectl set image deployment/affiliate-api \
    api=gcr.io/affiliate-platform-prod/affiliate-api:v1.2.3

# 3. Monitor rollout
kubectl rollout status deployment/affiliate-api --timeout=300s

# 4. Verify deployment
kubectl get pods -l app=affiliate-api
curl -f https://api.affiliate-platform.com/health

# 5. Run post-deployment tests
./scripts/post-deployment-tests.sh
```

**Rollback Procedure**:
```bash
# 1. Rollback to previous version
kubectl rollout undo deployment/affiliate-api

# 2. Verify rollback
kubectl rollout status deployment/affiliate-api
curl -f https://api.affiliate-platform.com/health

# 3. Notify team of rollback
echo "Deployment rolled back at $(date)" | \
    curl -X POST -H 'Content-type: application/json' \
    --data-binary @- \
    "$SLACK_WEBHOOK_URL"
```

### 5.2 Emergency Hotfix Deployment

**When to Use**: Critical security fixes, data corruption fixes

**Procedure**:
```bash
# 1. Build and push hotfix image
docker build -t gcr.io/affiliate-platform-prod/affiliate-api:hotfix-$(date +%Y%m%d-%H%M) .
docker push gcr.io/affiliate-platform-prod/affiliate-api:hotfix-$(date +%Y%m%d-%H%M)

# 2. Deploy immediately (skip staging)
kubectl set image deployment/affiliate-api \
    api=gcr.io/affiliate-platform-prod/affiliate-api:hotfix-$(date +%Y%m%d-%H%M)

# 3. Monitor closely
kubectl logs -f deployment/affiliate-api

# 4. Verify fix
# Run specific tests for the hotfix

# 5. Document emergency deployment
echo "Emergency hotfix deployed: $(date)" >> /var/log/emergency-deployments.log
```

## 6. Database Operations

### 6.1 Database Maintenance

**Weekly Maintenance**:
```bash
# 1. Update table statistics
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c 'ANALYZE;'
"

# 2. Vacuum tables
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c 'VACUUM (VERBOSE, ANALYZE);'
"

# 3. Check for bloated tables
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT schemaname, tablename, 
               pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
        FROM pg_tables 
        WHERE schemaname = 'public' 
        ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC 
        LIMIT 10;
    \"
"
```

**Performance Monitoring**:
```bash
# 1. Check slow queries
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT query, mean_time, calls, total_time
        FROM pg_stat_statements 
        ORDER BY mean_time DESC 
        LIMIT 10;
    \"
"

# 2. Check index usage
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT schemaname, tablename, indexname, idx_scan, idx_tup_read, idx_tup_fetch
        FROM pg_stat_user_indexes 
        WHERE idx_scan = 0;
    \"
"

# 3. Check connection statistics
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT state, count(*) 
        FROM pg_stat_activity 
        GROUP BY state;
    \"
"
```

### 6.2 Database Backup and Recovery

**Manual Backup**:
```bash
# 1. Create manual backup
gcloud sql backups create \
    --instance=affiliate-prod-db \
    --description="Manual backup $(date)"

# 2. Verify backup
gcloud sql backups list --instance=affiliate-prod-db --limit=5
```

**Point-in-Time Recovery**:
```bash
# 1. Create recovery instance
RECOVERY_TIME="2025-08-05 14:30:00"
gcloud sql instances create affiliate-recovery-temp \
    --source-instance=affiliate-prod-db \
    --source-instance-region=us-central1 \
    --point-in-time="$RECOVERY_TIME"

# 2. Validate recovered data
gcloud sql connect affiliate-recovery-temp --user=postgres

# 3. Promote recovery instance (after validation)
# This is a destructive operation - get approval first
# gcloud sql instances promote-replica affiliate-recovery-temp
```

## 7. Monitoring and Alerting

### 7.1 Key Metrics to Monitor

**Application Metrics**:
```bash
# Check current metrics
curl -s http://localhost:8080/metrics | grep -E "(http_requests_total|http_request_duration)"

# Key metrics to watch:
# - http_requests_total (request rate)
# - http_request_duration_seconds (response time)
# - database_connections_active (DB connections)
# - redis_commands_processed_total (Redis usage)
```

**Infrastructure Metrics**:
```bash
# Pod resource usage
kubectl top pods -l app=affiliate-api

# Node resource usage
kubectl top nodes

# Database metrics
gcloud sql instances describe affiliate-prod-db --format="table(
    name,
    state,
    settings.tier,
    settings.dataDiskSizeGb,
    settings.dataDiskType
)"
```

### 7.2 Alert Response Procedures

**High Error Rate Alert**:
```bash
# 1. Check error logs
kubectl logs -l app=affiliate-api --since=10m | grep -i error

# 2. Check error patterns
kubectl logs -l app=affiliate-api --since=1h | grep -i error | sort | uniq -c | sort -nr

# 3. If errors are from external service:
# - Check external service status
# - Enable circuit breaker if available
# - Consider temporary degraded mode
```

**High Response Time Alert**:
```bash
# 1. Check current response times
curl -w "%{time_total}\n" -o /dev/null -s https://api.affiliate-platform.com/health

# 2. Check database performance
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT query, mean_time 
        FROM pg_stat_statements 
        WHERE mean_time > 1000 
        ORDER BY mean_time DESC 
        LIMIT 5;
    \"
"

# 3. Scale if needed
kubectl scale deployment affiliate-api --replicas=5
```

**Database Connection Alert**:
```bash
# 1. Check connection count
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT count(*) as connections, 
               (SELECT setting FROM pg_settings WHERE name='max_connections') as max_connections;
    \"
"

# 2. Kill idle connections if needed
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT pg_terminate_backend(pid)
        FROM pg_stat_activity
        WHERE state = 'idle'
        AND state_change < now() - interval '30 minutes';
    \"
"
```

## 8. Security Operations

### 8.1 Security Incident Response

**Suspected Security Breach**:
```bash
# 1. Immediately isolate affected systems
kubectl scale deployment affiliate-api --replicas=0

# 2. Preserve evidence
kubectl logs -l app=affiliate-api --since=24h > security-incident-logs-$(date +%Y%m%d).txt

# 3. Check for suspicious activity
kubectl logs -l app=affiliate-api | grep -E "(failed|unauthorized|suspicious)"

# 4. Notify security team
# Send to: security@company.com
# Include: incident details, affected systems, initial assessment

# 5. Follow security incident response plan
# See Security & Compliance Guide for detailed procedures
```

**Certificate Expiry**:
```bash
# 1. Check certificate status
echo | openssl s_client -servername api.affiliate-platform.com -connect api.affiliate-platform.com:443 2>/dev/null | openssl x509 -noout -dates

# 2. Renew certificate (if using cert-manager)
kubectl delete certificate affiliate-platform-tls
kubectl apply -f k8s/certificates/

# 3. Verify new certificate
curl -vI https://api.affiliate-platform.com 2>&1 | grep -E "(expire|valid)"
```

### 8.2 Access Control Operations

**User Access Issues**:
```bash
# 1. Check user authentication
kubectl logs -l app=affiliate-api | grep "user_id:$USER_ID" | grep -i auth

# 2. Verify JWT token
# Use JWT debugger or:
echo "$JWT_TOKEN" | cut -d. -f2 | base64 -d | jq .

# 3. Check user permissions
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    psql \$DATABASE_URL -c \"
        SELECT id, email, role, organization_id 
        FROM profiles 
        WHERE id = '$USER_ID';
    \"
"
```

**API Key Issues**:
```bash
# 1. Check API key usage
kubectl logs -l app=affiliate-api | grep "api_key" | grep "$API_KEY"

# 2. Rotate compromised API key
kubectl create secret generic api-keys \
    --from-literal=stripe-key="$NEW_STRIPE_KEY" \
    --from-literal=everflow-key="$NEW_EVERFLOW_KEY" \
    --dry-run=client -o yaml | kubectl apply -f -

# 3. Restart pods to pick up new keys
kubectl rollout restart deployment/affiliate-api
```

## 9. Backup and Recovery Operations

### 9.1 Backup Verification

**Daily Backup Check**:
```bash
# 1. Verify latest backup exists
LATEST_BACKUP=$(gcloud sql backups list --instance=affiliate-prod-db --limit=1 --format="value(startTime)")
echo "Latest backup: $LATEST_BACKUP"

# 2. Check backup age
BACKUP_AGE=$(( ($(date +%s) - $(date -d "$LATEST_BACKUP" +%s)) / 3600 ))
if [ $BACKUP_AGE -gt 25 ]; then
    echo "WARNING: Backup is $BACKUP_AGE hours old"
else
    echo "Backup age OK: $BACKUP_AGE hours"
fi

# 3. Test backup integrity (monthly)
# Create test restore instance and validate data
```

### 9.2 Recovery Operations

**Application Recovery**:
```bash
# 1. Restore from Git
git checkout $LAST_KNOWN_GOOD_COMMIT
docker build -t recovery-image .
kubectl set image deployment/affiliate-api api=recovery-image

# 2. Restore configuration
kubectl apply -f k8s/overlays/production/

# 3. Verify recovery
kubectl rollout status deployment/affiliate-api
curl -f https://api.affiliate-platform.com/health
```

**Data Recovery**:
```bash
# 1. Stop application writes
kubectl scale deployment affiliate-api --replicas=0

# 2. Restore database
# Follow procedures in BCP/DR Plan

# 3. Restart application
kubectl scale deployment affiliate-api --replicas=3

# 4. Verify data integrity
# Run data validation scripts
```

## 10. Performance Optimization

### 10.1 Resource Scaling

**Horizontal Scaling**:
```bash
# Scale up during high traffic
kubectl scale deployment affiliate-api --replicas=8

# Scale down during low traffic
kubectl scale deployment affiliate-api --replicas=3

# Auto-scaling (if HPA is configured)
kubectl autoscale deployment affiliate-api --cpu-percent=70 --min=3 --max=10
```

**Vertical Scaling**:
```bash
# Increase resource limits
kubectl patch deployment affiliate-api -p '{
    "spec": {
        "template": {
            "spec": {
                "containers": [{
                    "name": "api",
                    "resources": {
                        "limits": {"memory": "1Gi", "cpu": "1000m"},
                        "requests": {"memory": "512Mi", "cpu": "500m"}
                    }
                }]
            }
        }
    }
}'
```

### 10.2 Cache Optimization

**Redis Cache Management**:
```bash
# Check cache hit rate
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    redis-cli -h \$REDIS_HOST info stats | grep keyspace_hits
"

# Clear specific cache patterns
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "
    redis-cli -h \$REDIS_HOST --scan --pattern 'cache:old:*' | xargs redis-cli -h \$REDIS_HOST del
"

# Warm up cache
curl -X POST http://localhost:8080/internal/cache/warm
```

---

## Appendix A: Quick Reference Commands

### Essential Commands
```bash
# Service health check
curl -f https://api.affiliate-platform.com/health

# Check pod status
kubectl get pods -l app=affiliate-api

# View recent logs
kubectl logs -l app=affiliate-api --tail=50

# Check resource usage
kubectl top pods -l app=affiliate-api

# Scale deployment
kubectl scale deployment affiliate-api --replicas=5

# Restart deployment
kubectl rollout restart deployment/affiliate-api

# Check database connections
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "psql \$DATABASE_URL -c 'SELECT count(*) FROM pg_stat_activity;'"

# Test Redis connectivity
kubectl exec -it deployment/affiliate-api -- /bin/sh -c "redis-cli -h \$REDIS_HOST ping"
```

### Emergency Commands
```bash
# Emergency scale down (if under attack)
kubectl scale deployment affiliate-api --replicas=0

# Emergency rollback
kubectl rollout undo deployment/affiliate-api

# Emergency database failover
gcloud sql instances promote-replica affiliate-prod-db-replica

# Update status page (service down)
curl -X PATCH "https://api.statuspage.io/v1/pages/$PAGE_ID/components/$COMPONENT_ID" \
    -H "Authorization: OAuth $STATUSPAGE_API_KEY" \
    -d '{"component": {"status": "major_outage"}}'
```

## Appendix B: Contact Information

### On-Call Escalation
1. **Primary On-Call**: +1-555-0199 (PagerDuty)
2. **Secondary On-Call**: +1-555-0198 (PagerDuty)
3. **SRE Lead**: Alex Rodriguez (+1-555-0103)
4. **Service Manager**: Mike Chen (+1-555-0102)

### External Vendors
- **Google Cloud Support**: +1-877-453-6021
- **Supabase Support**: support@supabase.com
- **Stripe Support**: +1-888-963-8331
- **Cloudflare Support**: +1-888-993-5273

---

**Document Classification**: Operational  
**Access Level**: On-Call Team, Operations Team  
**Review Frequency**: Monthly  
**Related Documents**: System Architecture, Monitoring Playbook, BCP/DR Plan