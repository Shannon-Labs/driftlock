# Driftlock Runbooks

This document provides operational runbooks for common tasks and troubleshooting procedures for the Driftlock platform.

## Table of Contents
- [Platform Overview](#platform-overview)
- [Deployment Runbooks](#deployment-runbooks)
- [Monitoring & Alerting](#monitoring--alerting)
- [Troubleshooting Scenarios](#troubleshooting-scenarios)
- [Incident Response](#incident-response)
- [Backup & Recovery](#backup--recovery)
- [Performance Tuning](#performance-tuning)

## Platform Overview

Driftlock is a production-ready anomaly detection platform that identifies deviations in telemetry data using compression-based anomaly detection (CBAD). The platform is designed for regulated environments requiring explainable and deterministic results.

### Architecture Overview
- **API Layer**: Go-based REST API handling requests and responses
- **Storage**: PostgreSQL for persistence of anomalies and configuration
- **CBAD Engine**: Rust-based compression engine performing anomaly detection
- **OTel Integration**: OpenTelemetry collector processor integration
- **Frontend**: Next.js UI for visualization and management

## Deployment Runbooks

### Production Deployment

#### Prerequisites
- Kubernetes cluster (v1.22+)
- Helm 3.0+
- PostgreSQL 13+ (external or via Helm chart)
- Load balancer for external access

#### Steps
1. **Clone the repository and navigate to the Helm chart:**
```bash
git clone https://github.com/shannon-labs/driftlock.git
cd driftlock/helm/driftlock
```

2. **Customize the values.yaml file for your environment**
```bash
cp values.yaml values.production.yaml
# Edit values.production.yaml with your specific configuration
```

3. **Create required secrets:**
```bash
# Generate JWT secret (at least 32 characters)
export JWT_SECRET=$(openssl rand -base64 32)
kubectl create secret generic driftlock-secrets \
  --from-literal=jwt-secret="$JWT_SECRET" \
  --from-literal=database-url="postgresql://user:pass@host:5432/dbname?sslmode=disable"
```

4. **Install the Helm chart:**
```bash
helm install driftlock . -f values.production.yaml -n driftlock --create-namespace
```

5. **Verify the deployment:**
```bash
kubectl get pods -n driftlock
kubectl get services -n driftlock
kubectl logs -l app=driftlock-api -n driftlock
```

6. **Run database migrations:**
```bash
kubectl exec -it deployment/driftlock-api -n driftlock -- \
  ./driftlock-api migrate
```

#### Post-Deployment Verification
- Verify all pods are running: `kubectl get pods -n driftlock`
- Check service is accessible: `kubectl get service driftlock-api -n driftlock`
- Test API health endpoint: `curl http://<service-ip>/healthz`
- Verify database connectivity in logs

### Upgrade Process

#### Pre-Upgrade Checklist
- [ ] Backup database
- [ ] Review release notes for breaking changes
- [ ] Test upgrade in staging environment
- [ ] Ensure sufficient cluster resources

#### Steps
1. **Scale down traffic to the service (if using ingress)**
2. **Export current configuration:**
```bash
kubectl get configmap driftlock-config -n driftlock -o yaml > backup-config.yaml
```

3. **Update Helm repository:**
```bash
helm repo update
```

4. **Perform the upgrade:**
```bash
helm upgrade driftlock driftlock/driftlock -f values.production.yaml -n driftlock
```

5. **Monitor the rollout:**
```bash
kubectl rollout status deployment/driftlock-api -n driftlock
kubectl get pods -w -n driftlock
```

6. **Verify functionality after upgrade:**
- Check health endpoints
- Verify metrics are being collected
- Validate that new features work as expected

## Monitoring & Alerting

### Key Metrics to Monitor

#### System Metrics
- `driftlock_goroutines` - Number of running goroutines
- `driftlock_memory_allocated_bytes` - Memory allocation
- `driftlock_db_connections_total` - Database connection pool usage
- `driftlock_http_requests_total` - Total HTTP requests by status code
- `driftlock_http_request_duration_seconds` - Request duration percentiles

#### Business Metrics
- `driftlock_anomalies_detected_total` - Total anomalies detected by stream type and severity
- `driftlock_anomaly_processing_duration_seconds` - Time to process anomalies
- `driftlock_cbad_events_per_second` - Processing rate of the CBAD engine
- `driftlock_cbad_compression_ratio` - Compression ratios for different data types

### Alerting Rules

#### High Priority Alerts
- **API Server Down**: `up{job="driftlock-api"} == 0`
- **High Error Rate**: `rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05`
- **High Latency**: `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1`
- **Database Connectivity**: `driftlock_db_connections_total{connection_state="failed"} > 0`

#### Medium Priority Alerts
- **High Anomaly Rate**: `rate(driftlock_anomalies_detected_total[5m]) > 100`
- **Memory Pressure**: `driftlock_memory_heap_bytes / (1024*1024*1024) > 0.8`
- **Connection Pool Exhaustion**: `driftlock_db_connections_total{connection_state="max"} > 0`

#### Low Priority Alerts
- **Low Processing Rate**: `driftlock_cbad_events_per_second < 10`
- **Configuration Change**: `sum(changes(driftlock_anomaly_ncd_threshold[5m])) > 0`

### Setting Up Monitoring Infrastructure

#### Prometheus Configuration
Add the following job to your Prometheus configuration:

```yaml
- job_name: 'driftlock-api'
  kubernetes_sd_configs:
  - role: pod
    namespaces:
      names:
      - driftlock
  relabel_configs:
  - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
    action: keep
    regex: true
  - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
    action: replace
    target_label: __metrics_path__
    regex: (.+)
  - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
    action: replace
    regex: ([^:]+)(?::\d+)?;(\d+)
    replacement: $1:$2
    target_label: __address__
```

#### Grafana Dashboard Setup
1. Import the dashboard from `monitoring/grafana/dashboard.json`
2. Configure Prometheus as the data source
3. Set up alerting rules based on the metrics mentioned above

## Troubleshooting Scenarios

### High CPU Usage

#### Symptoms
- CPU usage consistently above 80%
- Slow response times
- Request timeouts

#### Diagnosis Steps
1. **Check for resource-intensive queries:**
```sql
-- Check for slow queries
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;
```

2. **Analyze application profiling:**
```bash
# Access profiling endpoint
kubectl port-forward service/driftlock-api 8080:80 -n driftlock
go tool pprof http://localhost:8080/debug/pprof/profile
```

3. **Check for memory leaks or high garbage collection:**
```bash
# Look for GC-related metrics
kubectl logs -l app=driftlock-api -n driftlock | grep GC
```

#### Resolution
- Optimize slow database queries
- Add indexes for frequently queried columns
- Tune garbage collection settings
- Scale application resources if needed

### Database Connection Issues

#### Symptoms
- API errors with database-related messages
- High response times
- Connection timeout errors

#### Diagnosis Steps
1. **Check database connectivity:**
```bash
# Test direct connection to database
kubectl run pg-test -it --rm --image=postgres:13 --restart=Never -- \
  psql -h <db-host> -U <db-user> -d <db-name>
```

2. **Check connection pool metrics:**
```bash
kubectl logs -l app=driftlock-api -n driftlock | grep -i pool
```

3. **Monitor database connections:**
```sql
-- Check active connections
SELECT count(*) FROM pg_stat_activity;
-- Check connection limits
SHOW max_connections;
```

#### Resolution
- Increase database connection limits
- Optimize connection pool configuration
- Check for connection leaks in application code
- Verify network connectivity

### Anomaly Detection Performance Issues

#### Symptoms
- Delays in anomaly detection
- High memory usage
- Processing lag

#### Diagnosis Steps
1. **Check CBAD processing metrics:**
```bash
# Check events per second metrics
kubectl logs -l app=driftlock-api -n driftlock | grep cbad
```

2. **Analyze resource usage:**
```bash
kubectl top pods -n driftlock
kubectl describe pod <pod-name> -n driftlock
```

3. **Check data ingestion rate:**
```bash
# If using OTel collector, check its metrics
kubectl logs -l app=otel-collector -n driftlock
```

#### Resolution
- Scale CBAD processing resources
- Optimize compression algorithm parameters
- Implement data sampling if ingestion rate is too high
- Consider horizontal scaling

## Incident Response

### Security Incident Response

#### Compromised API Keys
1. **Immediate Actions:**
   - Rotate the compromised API key immediately
   - Review logs for unauthorized access
   - Revoke access for the compromised key

2. **Investigation:**
   - Check audit logs for suspicious activities
   - Identify what data was accessed
   - Determine scope of potential data exposure

3. **Recovery:**
   - Issue new API keys
   - Update all services with new keys
   - Implement additional security measures if needed

#### Data Breach
1. **Immediate Actions:**
   - Isolate affected systems
   - Preserve evidence
   - Stop data exfiltration

2. **Assessment:**
   - Determine scope of breach
   - Identify compromised data
   - Assess potential impact

3. **Notification:**
   - Follow internal incident response procedures
   - Notify appropriate authorities if required
   - Communicate with affected parties

### Service Outage Response

#### API Service Unavailable
1. **Verify the issue:**
   - Check if the pods are running: `kubectl get pods -n driftlock`
   - Check service status: `kubectl get services -n driftlock`
   - Check logs: `kubectl logs deployment/driftlock-api -n driftlock`

2. **Common fixes:**
   - Restart pods: `kubectl rollout restart deployment/driftlock-api -n driftlock`
   - Check resource limits
   - Verify configuration

3. **Escalate if needed:**
   - If issue persists, check infrastructure
   - Review recent changes
   - Engage platform team if needed

## Backup & Recovery

### Database Backup

#### Automated Backups
The production setup should include automated database backups:

```bash
# Example backup script
pg_dump -h <db-host> -U <db-user> -d <db-name> | gzip > driftlock-backup-$(date +%Y%m%d-%H%M%S).sql.gz
```

#### Backup Schedule
- **Daily full backups** at 2 AM
- **Hourly incremental backups** during business hours
- **Retention**: 30 days for daily, 7 days for hourly

### Recovery Process

#### Full Recovery
1. **Stop the Driftlock services:**
```bash
kubectl scale deployment/driftlock-api --replicas=0 -n driftlock
```

2. **Restore the database:**
```bash
gunzip -c backup-file.sql.gz | psql -h <db-host> -U <db-user> -d <db-name>
```

3. **Start the services:**
```bash
kubectl scale deployment/driftlock-api --replicas=3 -n driftlock
```

#### Point-in-Time Recovery
For PostgreSQL with WAL archiving:
```bash
# Restore to a specific point in time
pg_restore -h <db-host> -U <db-user> -d <db-name> --to-timestamp="2023-06-01 10:30:00"
```

## Performance Tuning

### Database Tuning

#### Connection Pool Settings
```yaml
# In ConfigMap or values file
DB_MAX_OPEN_CONNS: "25"
DB_MAX_IDLE_CONNS: "10"
DB_CONN_MAX_LIFETIME: "1h"
```

#### PostgreSQL Configuration
```sql
-- Recommended settings for high-load scenarios
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET work_mem = '4MB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET max_connections = '100';
```

### Application Tuning

#### Memory and GC Optimization
```bash
# Environment variables for Go runtime
GOGC=50          # Garbage collection target percentage
GOMAXPROCS=4     # Number of operating system threads
GOMEMLIMIT=1GiB  # Memory limit for garbage collection
```

#### Rate Limiting Configuration
```yaml
# In ConfigMap
GLOBAL_RATE_LIMIT_RPS: "1000"
GLOBAL_RATE_LIMIT_BURST: "2000"
```

### Horizontal Scaling

#### Auto-scaling Configuration
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: driftlock-api
  namespace: driftlock
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: driftlock-api
  minReplicas: 3
  maxReplicas: 10
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
```

## Maintenance Procedures

### Regular Maintenance Tasks

#### Weekly
- Review logs for unusual patterns
- Check backup job completion
- Verify alerting rules are functioning

#### Monthly
- Performance review and optimization
- Security patching
- Configuration review and updates
- Capacity planning

#### Quarterly
- Disaster recovery drill
- Complete system audit
- Update documentation
- Review and update runbooks

### Health Checks

#### Automated Health Checks
```bash
# API health check
curl -f http://<api-endpoint>/healthz

# Database connectivity check
curl -f http://<api-endpoint>/healthz/db

# Configuration check
curl -f http://<api-endpoint>/healthz/config
```

#### Manual Health Verification
- Verify all pods are running and healthy
- Check metrics are being collected
- Test API endpoints manually
- Validate recent anomaly detections