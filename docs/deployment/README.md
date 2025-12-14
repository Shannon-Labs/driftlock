# Driftlock Deployment Guides

This directory contains deployment guides for Driftlock on various cloud platforms and environments.

## Ingestion Guides

- **[Kafka Integration](./KAFKA_INTEGRATION.md)** - High-throughput streaming via Kafka consumer
- **[OTLP Ingestion](./OTLP_INGESTION.md)** - Native OpenTelemetry Protocol gRPC server
- **[Kafka Setup](./KAFKA_SETUP.md)** - Local Kafka/Zookeeper setup

## Build Features

```bash
# HTTP only (default)
cargo build -p driftlock-api --release

# With Kafka consumer
cargo build -p driftlock-api --features kafka --release

# With OTLP gRPC server
cargo build -p driftlock-api --features otlp --release

# All features
cargo build -p driftlock-api --features kafka,otlp,webhooks --release
```

## Overview

Driftlock can be deployed in multiple environments:

- **Kubernetes**: Production deployments using Helm charts
- **Docker Compose**: Development and demo environments
- **Cloud Platforms**: AWS, Azure, GCP, and other cloud providers

## Quick Start

### Docker Compose

```bash
# Clone the repository
git clone https://github.com/Shannon-Labs/driftlock.git
cd driftlock

# Start all services
docker-compose up -d
```

### Kubernetes with Helm

```bash
# Add the Driftlock Helm repository
helm repo add driftlock https://charts.driftlock.com

# Install for development
helm install driftlock/driftlock --create-namespace driftlock --set api.replicaCount=1

# Install for production
helm install driftlock/driftlock \
  --create-namespace driftlock \
  --set api.replicaCount=3 \
  --set kafka.enabled=true \
  --set clickhouse.enabled=true \
  --set prometheus.enabled=true \
  --set grafana.enabled=true
```

## Cloud Platform Guides

### AWS

#### Prerequisites

- AWS CLI installed and configured
- kubectl configured to access EKS cluster
- Helm 3.x installed

#### Deployment

1. **Create EKS Cluster**:
   ```bash
   eksctl create cluster --name driftlock --region us-west-2 --nodegroup-name standard --node-type m5.large --nodes 3
   ```

2. **Configure kubectl**:
   ```bash
   aws eks update-kubeconfig --name driftlock --region us-west-2
   ```

3. **Deploy Driftlock**:
   ```bash
   helm install driftlock/driftlock \
     --namespace driftlock \
     --set api.replicaCount=3 \
     --set kafka.enabled=true \
     --set clickhouse.enabled=true \
     --set prometheus.enabled=true \
     --set grafana.enabled=true \
     --set database.external.enabled=true \
     --set database.external.host=your-rds-endpoint.rds.amazonaws.com \
     --set database.external.database=driftlock \
     --set database.external.username=driftlock_user \
     --set database.external.password=your-secure-password
   ```

4. **Configure Ingress**:
   ```bash
   kubectl apply -f - <<EOF
   apiVersion: networking.k8s.io/v1
   kind: Ingress
   metadata:
     name: driftlock-ingress
     namespace: driftlock
     annotations:
       kubernetes.io/ingress.class: nginx
       cert-manager.io/cluster-issuer: letsencrypt-prod
   spec:
     tls:
       - hosts:
         - api.driftlock.com
     rules:
       - host: api.driftlock.com
         http:
           paths:
             - path: /
               pathType: Prefix
           backend:
             service:
               name: driftlock-api
               port:
                 number: 8080
   EOF
   ```

### Azure

#### Prerequisites

- Azure CLI installed and configured
- kubectl configured to access AKS cluster
- Helm 3.x installed

#### Deployment

1. **Create AKS Cluster**:
   ```bash
   az aks create --resource-group driftlock-rg --name driftlock --node-count 3 --node-vm-size Standard_D2s_v3 --enable-addons monitoring
   ```

2. **Configure kubectl**:
   ```bash
   az aks get-credentials --resource-group driftlock-rg --name driftlock --file ~/.kube/config
   ```

3. **Deploy Driftlock**:
   ```bash
   helm install driftlock/driftlock \
     --namespace driftlock \
     --set api.replicaCount=3 \
     --set kafka.enabled=true \
     --set clickhouse.enabled=true \
     --set prometheus.enabled=true \
     --set grafana.enabled=true \
     --set database.external.enabled=true \
     --set database.external.host=your-driftlock-db.postgres.database.azure.com \
     --set database.external.database=driftlock \
     --set database.external.username=driftlock_user@your-driftlock-db \
     --set database.external.password=your-secure-password
   ```

4. **Configure Ingress**:
   ```bash
   kubectl apply -f - <<EOF
   apiVersion: networking.k8s.io/v1
   kind: Ingress
   metadata:
     name: driftlock-ingress
     namespace: driftlock
     annotations:
       kubernetes.io/ingress.class: nginx
       cert-manager.io/cluster-issuer: letsencrypt-prod
   spec:
     tls:
       - hosts:
         - api.driftlock.com
     rules:
       - host: api.driftlock.com
         http:
           paths:
             - path: /
               pathType: Prefix
           backend:
             service:
               name: driftlock-api
               port:
                 number: 8080
   EOF
   ```

### GCP

#### Prerequisites

- gcloud CLI installed and configured
- kubectl configured to access GKE cluster
- Helm 3.x installed

#### Deployment

1. **Create GKE Cluster**:
   ```bash
   gcloud container clusters create driftlock --num-nodes 3 --machine-type e2-standard-4 --zone us-central1-a
   ```

2. **Configure kubectl**:
   ```bash
   gcloud container clusters get-credentials driftlock --zone us-central1-a
   ```

3. **Deploy Driftlock**:
   ```bash
   helm install driftlock/driftlock \
     --namespace driftlock \
     --set api.replicaCount=3 \
     --set kafka.enabled=true \
     --set clickhouse.enabled=true \
     --set prometheus.enabled=true \
     --set grafana.enabled=true \
     --set database.external.enabled=true \
     --set database.external.host=your-driftlock-db.postgres.database.azure.com \
     --set database.external.database=driftlock \
     --set database.external.username=driftlock_user@your-driftlock-db \
     --set database.external.password=your-secure-password
   ```

## Configuration

### Environment Variables

| Variable | Description | Default |
|-----------|-------------|---------|
| API_REPLICAS | Number of API replicas | 3 |
| KAFKA_ENABLED | Enable Kafka streaming | true |
| CLICKHOUSE_ENABLED | Enable ClickHouse analytics | true |
| PROMETHEUS_ENABLED | Enable Prometheus monitoring | true |
| GRAFANA_ENABLED | Enable Grafana dashboards | true |
| DATABASE_EXTERNAL_ENABLED | Use external database | false |
| DATABASE_HOST | Database host | localhost |
| DATABASE_PORT | Database port | 5432 |
| DATABASE_NAME | Database name | driftlock |
| DATABASE_USER | Database user | postgres |
| DATABASE_PASSWORD | Database password | password |

### Resource Requirements

| Component | Minimum | Recommended |
|-----------|----------|-------------|
| API Server | 1 CPU, 512MB RAM | 2 CPU, 1GB RAM |
| Kafka | 1 CPU, 1GB RAM | 3 CPU, 3GB RAM |
| ClickHouse | 1 CPU, 2GB RAM | 2 CPU, 4GB RAM |
| PostgreSQL | 1 CPU, 1GB RAM | 2 CPU, 2GB RAM |
| Redis | 0.5 CPU, 256MB RAM | 1 CPU, 512MB RAM |
| Prometheus | 1 CPU, 1GB RAM | 2 CPU, 2GB RAM |
| Grafana | 0.5 CPU, 512MB RAM | 1 CPU, 1GB RAM |

### Scaling Guidelines

#### Horizontal Scaling

- **API Server**: Use HPA based on CPU and memory usage
- **Kafka**: Increase partition count and broker instances
- **ClickHouse**: Add more nodes to the cluster
- **PostgreSQL**: Use read replicas and connection pooling

#### Vertical Scaling

- **API Server**: Increase CPU and memory limits
- **Kafka**: Increase broker memory and disk I/O
- **ClickHouse**: Increase node memory and CPU

### Monitoring

#### Prometheus Metrics

Key metrics to monitor:

- **System**: CPU, memory, disk, network
- **Application**: Request rate, response time, error rate
- **Anomaly Detection**: Detection rate, types, confidence scores
- **Tenant**: Resource usage per tenant

#### Grafana Dashboards

- **Overview**: System health and performance
- **Anomaly Detection**: Anomaly analysis and trends
- **Tenant**: Per-tenant resource usage and quotas
- **Performance**: System performance metrics

### Security

#### Network Security

- Use network policies to restrict traffic between services
- Enable TLS encryption for all external communication
- Implement rate limiting at the ingress level

#### Authentication

- Use API key authentication for all API access
- Rotate API keys regularly
- Implement proper secret management

### Backup and Recovery

#### Database

- Enable automated backups for PostgreSQL
- Configure point-in-time recovery for ClickHouse
- Test backup and restore procedures regularly

#### Application State

- Use Redis for session storage and caching
- Implement stateful processing with Kafka offsets
- Configure proper graceful shutdown

## Troubleshooting

### Common Issues

1. **Services Not Starting**: Check network connectivity and dependencies
2. **High Memory Usage**: Check resource limits and leaks
3. **Database Connection Issues**: Verify connection strings and network policies
4. **Kafka Not Processing**: Check broker health and topic configuration
5. **No Anomalies Detected**: Adjust CBAD thresholds and data quality

### Debugging

#### API Server

```bash
# Check API logs
docker-compose logs api

# Check API health
curl http://localhost:8080/healthz

# Debug with verbose logging
LOG_LEVEL=debug docker-compose up api
```

#### Kafka

```bash
# Check Kafka logs
docker-compose logs kafka

# List topics
docker-compose exec kafka kafka-topics.sh --bootstrap-server kafka:9092 --list

# Consume messages for debugging
docker-compose exec kafka kafka-console-consumer.sh --bootstrap-server kafka:9092 --topic driftlock-events --from-beginning
```

#### ClickHouse

```bash
# Check ClickHouse logs
docker-compose logs clickhouse

# Connect to ClickHouse CLI
docker-compose exec clickhouse clickhouse-client

# Check table structure
SHOW CREATE TABLE driftlock.anomalies

# Query anomalies
SELECT * FROM driftlock.anomalies WHERE timestamp >= now() - INTERVAL 1 DAY LIMIT 10
```

## Support

For support and questions:

- Documentation: https://docs.driftlock.com
- Email: support@driftlock.com
- Community: https://community.driftlock.com
- Status Page: https://status.driftlock.com
