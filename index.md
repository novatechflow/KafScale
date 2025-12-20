---
layout: default
title: Kafscale — Stateless Kafka on S3
description: Kafka-compatible streaming with stateless brokers, S3-native storage, and Kubernetes-first operations.
---

<section class="hero">
  <p class="eyebrow">Kafka-compatible streaming. Stateless brokers. S3 storage. Self-hosted.</p>
  <h1>Stateless Kafka on S3, compatible with your clients.</h1>
  <p>Run Kafka APIs without stateful disks. Kafscale stores log segments in S3, keeps brokers ephemeral, and uses etcd for metadata so you scale fast and recover cleanly.</p>
  <div class="badge-row">
    <img alt="GitHub stars" src="https://img.shields.io/github/stars/novatechflow/kafscale?style=flat" />
    <img alt="License" src="https://img.shields.io/badge/license-Apache%202.0-blue" />
  </div>
  <div class="hero-actions">
    <a class="button" href="/quickstart/">Get started</a>
    <a class="button secondary" href="https://github.com/novatechflow/kafscale" target="_blank" rel="noreferrer">View on GitHub</a>
  </div>
</section>

<section class="section">
  <h2>Why teams adopt Kafscale</h2>
  <div class="grid">
    <div class="card">
      <h3>Stateless brokers</h3>
      <p>Spin brokers up and down without disk shuffles. S3 is the source of truth.</p>
    </div>
    <div class="card">
      <h3>S3-native durability</h3>
      <p>Immutable segments, predictable retention, and lifecycle policies in object storage.</p>
    </div>
    <div class="card">
      <h3>etcd metadata</h3>
      <p>Offsets, topics, and group state live in etcd for fast consensus and control.</p>
    </div>
    <div class="card">
      <h3>Kubernetes operator</h3>
      <p>CRDs manage clusters, topics, snapshots, and UI access with clean drift control.</p>
    </div>
  </div>
</section>

<section class="section">
  <h2>Quickstart in minutes</h2>
  <div class="grid">
    <div class="card">
      <h3>1. Install operator</h3>
      <pre class="code-block"><code>helm upgrade --install kafscale deploy/helm/kafscale \
  --namespace kafscale --create-namespace \
  --set operator.etcdEndpoints={} \
  --set operator.image.tag=latest \
  --set console.image.tag=latest</code></pre>
    </div>
    <div class="card">
      <h3>2. Create your first topic</h3>
      <pre class="code-block"><code>kubectl apply -n kafscale -f - &lt;&lt;'EOF'
apiVersion: kafscale.novatechflow.io/v1alpha1
kind: KafscaleTopic
metadata:
  name: orders
spec:
  clusterRef: demo
  partitions: 3
EOF</code></pre>
    </div>
    <div class="card">
      <h3>3. Produce and consume</h3>
      <pre class="code-block"><code>kafka-console-producer --bootstrap-server 127.0.0.1:9092 --topic orders
kafka-console-consumer --bootstrap-server 127.0.0.1:9092 --topic orders --from-beginning</code></pre>
    </div>
  </div>
  <div class="hero-actions">
    <a class="button secondary" href="/quickstart/">Read the full quickstart</a>
  </div>
</section>

<section class="section">
  <h2>Operations you can trust</h2>
  <div class="grid">
    <div class="card">
      <h3>Prometheus metrics</h3>
      <p>Track S3 health state, broker throughput, and operator snapshots.</p>
    </div>
    <div class="card">
      <h3>Scale instantly</h3>
      <p>Stateless brokers pair with HPA and avoid disk rebalancing delays.</p>
    </div>
    <div class="card">
      <h3>Clear failure modes</h3>
      <p>Startup gating, S3 health states, and etcd snapshot checks keep you safe.</p>
    </div>
  </div>
</section>

<section class="section">
  <h2>Deep technical reference</h2>
  <div class="grid">
    <div class="card">
      <h3>Protocol coverage</h3>
      <p>Support for Produce, Fetch, Metadata, consumer groups, and admin APIs.</p>
      <a class="button secondary" href="/api/">View API docs</a>
    </div>
    <div class="card">
      <h3>Storage format</h3>
      <p>Segment and index file layout plus S3 key structure for lifecycle design.</p>
      <a class="button secondary" href="/storage-format/">Explore storage format</a>
    </div>
    <div class="card">
      <h3>Security posture</h3>
      <p>Current TLS and auth status, plus the roadmap for SASL and ACLs.</p>
      <a class="button secondary" href="/security/">Read security overview</a>
    </div>
  </div>
</section>
