---
layout: doc
title: Runtime Settings
description: Broker, S3, etcd, consumer group, and operator runtime settings for KafScale.
permalink: /configuration/
nav_title: Runtime Settings
nav_order: 1
nav_group: References
---

# Runtime Settings

All configuration is done via environment variables. Set these in your Helm values, KafScaleCluster CRD, or container spec.

---

## Broker configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFSCALE_BROKER_ID` | auto | Node ID for group membership; auto-assigned if empty |
| `KAFSCALE_SEGMENT_BYTES` | `4194304` | Segment flush threshold in bytes (4MB) |
| `KAFSCALE_FLUSH_INTERVAL_MS` | `500` | Maximum time before flushing buffer to S3 |
| `KAFSCALE_CACHE_BYTES` | `1073741824` | Hot segment cache size in bytes (1GB) |
| `KAFSCALE_INDEX_CACHE_BYTES` | `104857600` | Index cache size in bytes (100MB) |
| `KAFSCALE_READAHEAD_SEGMENTS` | `2` | Number of segments to prefetch |
| `KAFSCALE_AUTO_CREATE_TOPICS` | `true` | Auto-create topics on first produce |
| `KAFSCALE_AUTO_CREATE_PARTITIONS` | `1` | Partition count for auto-created topics |
| `KAFSCALE_THROUGHPUT_WINDOW_SEC` | `60` | Window for throughput metrics calculation |
| `KAFSCALE_STARTUP_TIMEOUT_SEC` | `30` | Broker startup timeout |
| `KAFSCALE_LOG_LEVEL` | `info` | Log level: debug, info, warn, error |

---

## S3 configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFSCALE_S3_BUCKET` | required | S3 bucket for segments and snapshots |
| `KAFSCALE_S3_REGION` | `us-east-1` | S3 region |
| `KAFSCALE_S3_ENDPOINT` | | S3 endpoint override (for MinIO, etc) |
| `KAFSCALE_S3_READ_BUCKET` | | Optional read replica bucket (CRR/MRAP) |
| `KAFSCALE_S3_READ_REGION` | | Optional read replica region |
| `KAFSCALE_S3_READ_ENDPOINT` | | Optional read replica endpoint override |
| `KAFSCALE_S3_PATH_STYLE` | `false` | Use path-style addressing instead of virtual-hosted |
| `KAFSCALE_S3_KMS_ARN` | | KMS key ARN for server-side encryption (SSE-KMS) |

Credentials (if not using IAM roles):

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFSCALE_S3_ACCESS_KEY` | | AWS access key ID |
| `KAFSCALE_S3_SECRET_KEY` | | AWS secret access key |
| `KAFSCALE_S3_SESSION_TOKEN` | | AWS session token (for temporary credentials) |

Health thresholds:

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFSCALE_S3_LATENCY_WARN_MS` | `500` | Latency threshold for degraded state |
| `KAFSCALE_S3_LATENCY_CRIT_MS` | `2000` | Latency threshold for unavailable state |
| `KAFSCALE_S3_ERROR_RATE_WARN` | `0.01` | Error rate threshold for degraded (1%) |
| `KAFSCALE_S3_ERROR_RATE_CRIT` | `0.05` | Error rate threshold for unavailable (5%) |
| `KAFSCALE_S3_HEALTH_WINDOW_SEC` | `60` | Health metric sampling window |

---

## etcd configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFSCALE_ETCD_ENDPOINTS` | required | Comma-separated etcd endpoints |
| `KAFSCALE_ETCD_USERNAME` | | etcd basic auth username |
| `KAFSCALE_ETCD_PASSWORD` | | etcd basic auth password |
| `KAFSCALE_ETCD_CERT_FILE` | | Path to client certificate |
| `KAFSCALE_ETCD_KEY_FILE` | | Path to client key |
| `KAFSCALE_ETCD_CA_FILE` | | Path to CA certificate |

Example endpoints:

```bash
# Single node (dev)
KAFSCALE_ETCD_ENDPOINTS=http://etcd.kafscale.svc:2379

# Cluster (prod)
KAFSCALE_ETCD_ENDPOINTS=http://etcd-0.etcd:2379,http://etcd-1.etcd:2379,http://etcd-2.etcd:2379
```

---

## Consumer group settings

Session timeout and heartbeat intervals are negotiated by Kafka clients following protocol defaults. KafScale respects the values sent by clients:

| Client setting | Typical default | Description |
|----------------|-----------------|-------------|
| `session.timeout.ms` | `45000` | Time before consumer is considered dead |
| `heartbeat.interval.ms` | `3000` | Heartbeat frequency |
| `max.poll.interval.ms` | `300000` | Max time between polls |

These are client-side settings, not broker configuration.

---

## Operator configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFSCALE_OPERATOR_ETCD_ENDPOINTS` | | External etcd endpoints; empty uses managed etcd |
| `KAFSCALE_OPERATOR_ETCD_SNAPSHOT_BUCKET` | | S3 bucket for etcd snapshots |
| `KAFSCALE_OPERATOR_ETCD_SNAPSHOT_PREFIX` | `etcd-snapshots/` | S3 key prefix for snapshots |
| `KAFSCALE_OPERATOR_ETCD_SNAPSHOT_SCHEDULE` | `0 */6 * * *` | Cron schedule for snapshots |
| `KAFSCALE_OPERATOR_ETCD_SNAPSHOT_S3_ENDPOINT` | | S3 endpoint override for snapshots |
| `KAFSCALE_OPERATOR_ETCD_SNAPSHOT_STALE_AFTER_SEC` | `3600` | Snapshot staleness alert threshold |
| `KAFSCALE_OPERATOR_RECONCILE_INTERVAL_SEC` | `30` | Reconciliation loop interval |
| `KAFSCALE_OPERATOR_LOG_LEVEL` | `info` | Operator log level |

---

## TLS configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFSCALE_TLS_ENABLED` | `false` | Enable TLS for client connections |
| `KAFSCALE_TLS_CERT_FILE` | | Path to server certificate |
| `KAFSCALE_TLS_KEY_FILE` | | Path to server private key |
| `KAFSCALE_TLS_CA_FILE` | | Path to CA certificate (for mTLS) |
| `KAFSCALE_TLS_CLIENT_AUTH` | `false` | Require client certificates |

See [Security](/security/) for TLS setup instructions.

---

## Example: minimal production config

```yaml
apiVersion: kafscale.io/v1alpha1
kind: KafScaleCluster
metadata:
  name: prod
  namespace: kafscale
spec:
  brokers:
    replicas: 3
    env:
      - name: KAFSCALE_SEGMENT_BYTES
        value: "16777216"
      - name: KAFSCALE_FLUSH_INTERVAL_MS
        value: "1000"
      - name: KAFSCALE_CACHE_BYTES
        value: "4294967296"
      - name: KAFSCALE_LOG_LEVEL
        value: "warn"
  s3:
    bucket: kafscale-prod
    region: us-east-1
    credentialsSecretRef: kafscale-s3
  etcd:
    endpoints:
      - http://etcd-0.etcd.kafscale.svc:2379
      - http://etcd-1.etcd.kafscale.svc:2379
      - http://etcd-2.etcd.kafscale.svc:2379
```

---

## Next steps

- [Installation](/installation/) for Helm and CRD setup
- [Operations](/operations/) for monitoring and scaling
- [Security](/security/) for TLS and authentication
