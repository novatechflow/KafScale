---
layout: doc
title: User Guide
description: How to interact with KafScale once it is deployed.
permalink: /user-guide/
nav_title: User Guide
nav_order: 4
---

<!--
Copyright 2025 Alexander Alten (novatechflow), NovaTechflow (novatechflow.com).
This project is supported and financed by Scalytics, Inc. (www.scalytics.io).

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# KafScale User Guide

KafScale is a Kafka-compatible, S3-backed message transport system. It keeps brokers stateless, stores data in S3, and relies on Kubernetes for scheduling and scaling. This guide summarizes how to interact with the platform once it is deployed.

## Concepts

- **Topics / Partitions**: match upstream Kafka semantics. All Kafka client libraries continue to work.
- **Brokers**: stateless pods accepting Kafka protocol traffic on port 9092 and metrics + gRPC control on 9093.
- **Metadata**: stored in etcd, encoded via protobufs (`kafscale.metadata.*`).
- **Storage**: message segments live in S3 buckets; brokers only keep in-memory caches.
- **Operator**: Kubernetes controller that provisions brokers, topics, and wiring based on CRDs.

## Before you start

The User Guide assumes you already have a cluster deployed. If you still need to deploy or configure the platform, use:

- [Quickstart](/quickstart/) for the shortest path to a working cluster
- [Installation](/installation/) for Helm values, CRDs, and environment setup

## Day-2 usage

Once the platform is up, day-2 operations typically include:

- Connecting existing Kafka clients to the broker service
- Monitoring broker health and metrics
- Planning scaling and maintenance workflows with the operator

For operational workflows, see [Operations](/operations/).

## Multi-Region S3 Reads (CRR)

If you run brokers in multiple regions, configure a read replica bucket per region so brokers read locally and fall back to the primary on CRR lag. Configure `spec.s3.readBucket`, `spec.s3.readRegion`, and `spec.s3.readEndpoint` in the cluster spec, or use the corresponding `KAFSCALE_S3_READ_*` environment variables.

For setup details, see [Operations](/operations/) and [Runtime Settings](/configuration/).

## Limits / Non-Goals

- No embedded stream processing features—pair KafScale with Flink, Wayang, Spark, etc.
- Transactions, idempotent producers, and log compaction are out of scope for the MVP.

For deeper architectural details or development guidance, read `kafscale-spec.md` and `docs/development.md`.
