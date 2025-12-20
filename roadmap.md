---
layout: doc
title: Roadmap
description: Completed milestones, current work, and what is planned next.
---

# Roadmap

## What’s done

- Core protocol parsing and metadata support
- Produce and fetch paths with S3-backed durability
- Consumer group coordination with offset and group persistence
- DescribeGroups/ListGroups ops visibility
- OffsetForLeaderEpoch consumer recovery
- DescribeConfigs/AlterConfigs ops tuning
- CreatePartitions/DeleteGroups ops APIs
- etcd topic/partition management
- Observability (structured logging, Grafana templates, Prometheus metrics)
- Kubernetes operator with managed etcd + snapshot backups
- End-to-end tests for broker durability and operator resilience
- Admin ops API e2e coverage
- Security review (TLS/auth)
- End-to-end tests multi-segment restart durability

## What’s in progress

- Performance benchmarks

## What’s planned

See the GitHub roadmap issues for near-term planning and design specs.

## How to request features

Open a GitHub issue with your use case, expected throughput, and operational constraints.
