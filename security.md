---
layout: doc
title: Security Overview
description: Current security posture, hardening guidance, and roadmap for KafScale.
permalink: /security/
---

# Security Overview

KafScale is a Kubernetes-native platform focused on Kafka protocol parity and operational stability. This document summarizes the current security posture and the boundaries of what is and is not supported in v1.

---

## Current security posture (v1)

| Area | Status |
|------|--------|
| Authentication | None at Kafka protocol layer; console supports basic auth |
| Authorization | None; admin APIs are unauthenticated |
| Transport | TLS optional, operator-configured |
| Secrets | S3 credentials via K8s secrets, not stored in etcd |
| Data at rest | Depends on S3/etcd provider encryption |
| Network | Assumes private network or cluster-level controls |

### Details

- **Authentication**: Brokers accept any client connection. The console UI supports basic auth via Helm values `console.auth.username` / `console.auth.password`.

- **Authorization**: None. Admin APIs are unauthenticated. The console UI is read-only by policy, not enforcement.

- **Transport security**: TLS is optional and must be enabled via environment variables:

| Variable | Description |
|----------|-------------|
| `KAFSCALE_TLS_ENABLED` | Enable TLS for broker connections |
| `KAFSCALE_TLS_CERT_FILE` | Path to server certificate |
| `KAFSCALE_TLS_KEY_FILE` | Path to server private key |
| `KAFSCALE_TLS_CA_FILE` | Path to CA certificate (for mTLS) |

- **Secrets handling**: S3 credentials are read from Kubernetes secrets and are never written to etcd or logged.

- **Data at rest**: Stored in S3 and etcd. Encryption depends on your provider configuration (S3 SSE, etcd encryption at rest).

---

## Operational guidance

### Network isolation

Deploy brokers and console behind private networking:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: kafscale-broker
  namespace: kafscale
spec:
  podSelector:
    matchLabels:
      app: kafscale-broker
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: allowed-namespace
      ports:
        - protocol: TCP
          port: 9092
```

### Enable TLS in production

```yaml
spec:
  brokers:
    env:
      - name: KAFSCALE_TLS_ENABLED
        value: "true"
      - name: KAFSCALE_TLS_CERT_FILE
        value: /etc/kafscale/tls/tls.crt
      - name: KAFSCALE_TLS_KEY_FILE
        value: /etc/kafscale/tls/tls.key
    volumeMounts:
      - name: tls
        mountPath: /etc/kafscale/tls
        readOnly: true
  volumes:
    - name: tls
      secret:
        secretName: kafscale-tls
```

### S3 IAM least privilege

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::kafscale-data",
        "arn:aws:s3:::kafscale-data/*"
      ]
    }
  ]
}
```

### etcd access control

- Use separate etcd credentials for KafScale
- Enable etcd authentication and RBAC
- Restrict etcd endpoints to broker and operator pods only

### Console protection

- Do not expose the console publicly without authentication
- Use ingress with authentication (OAuth proxy, basic auth)
- Consider disabling the console in production if not needed

---

## Known gaps

| Gap | Impact | Mitigation |
|-----|--------|------------|
| No SASL authentication | Any client can connect | Network isolation, VPN |
| No mTLS for clients | Cannot verify client identity | Network policies |
| No ACLs/RBAC | Cannot restrict topic access | Single-tenant deployments |
| No multi-tenant isolation | All clients see all topics | Separate clusters per tenant |
| Admin APIs unauthenticated | Anyone with network access can modify | Network isolation |

---

## Roadmap

| Feature | Target | Status |
|---------|--------|--------|
| TLS enabled by default | v1.1 | Planned |
| SASL/PLAIN authentication | v1.2 | Planned |
| SASL/SCRAM authentication | v1.2 | Planned |
| ACL layer for topics | v1.3 | Under design |
| mTLS for clients | v1.3 | Under design |
| Audit logging | v1.4 | Proposed |

---

## Reporting security issues

If you discover a security vulnerability, please report it responsibly:

1. **Do not** open a public GitHub issue
2. Email security@novatechflow.com with details
3. Include steps to reproduce if possible
4. We will respond within 48 hours

See [SECURITY.md](https://github.com/novatechflow/kafscale/blob/main/SECURITY.md) for full details.