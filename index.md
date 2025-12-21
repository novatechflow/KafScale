---
layout: default
title: KafScale - Stateless Kafka on S3
description: Kafka-compatible streaming with stateless brokers, S3-native storage, and Kubernetes-first operations. Apache 2.0 licensed.
---

<section class="hero">
  <p class="eyebrow">Apache 2.0 licensed. No vendor lock-in. Self-hosted.</p>
  <h1>Stateless Kafka on S3, compatible with your clients.</h1>
  <p>Run Kafka APIs without stateful disks. KafScale stores segments in S3, keeps brokers ephemeral, and uses etcd for metadata. Scale fast, recover cleanly, pay only for storage.</p>
  <div class="badge-row">
    <img alt="GitHub stars" src="https://img.shields.io/github/stars/novatechflow/kafscale?style=flat" />
    <img alt="License" src="https://img.shields.io/badge/license-Apache%202.0-blue" />
    <img alt="Go version" src="https://img.shields.io/github/go-mod/go-version/novatechflow/kafscale" />
  </div>
  <div class="hero-actions">
    <a class="button" href="/quickstart/">Get started</a>
    <a class="button secondary" href="https://github.com/novatechflow/kafscale" target="_blank" rel="noreferrer">View on GitHub</a>
  </div>
</section>

<section class="section">
  <h2>What teams are saying</h2>
  <div class="grid">
    <div class="card">
      <p>"After WarpStream got acquired, KafScale became our go-to. Better S3 integration, lower latency than we expected, fully scalable, and minimal ops burden."</p>
      <p><strong>— Platform team, Series B fintech</strong></p>
    </div>
    <div class="card">
      <p>"We moved 50 topics off Kafka in a weekend. No more disk alerts, no more partition rebalancing. Our on-call rotation got a lot quieter."</p>
      <p><strong>— SRE lead, e-commerce platform</strong></p>
    </div>
    <div class="card">
      <p>"The Apache 2.0 license was the deciding factor. We can't build on BSL projects, and we won't depend on a vendor's control plane."</p>
      <p><strong>— CTO, healthcare data startup</strong></p>
    </div>
  </div>
</section>

<section class="section">
  <h2>Why teams adopt KafScale</h2>
  <div class="grid">
    <div class="card">
      <h3>Stateless brokers</h3>
      <p>Spin brokers up and down without disk shuffles. S3 is the source of truth. No partition rebalancing, ever.</p>
    </div>
    <div class="card">
      <h3>S3-native durability</h3>
      <p>11 nines of durability. Immutable segments, lifecycle-based retention, predictable costs.</p>
    </div>
    <div class="card">
      <h3>Kubernetes operator</h3>
      <p>CRDs for clusters, topics, and snapshots. HPA-ready scaling. GitOps-friendly.</p>
    </div>
    <div class="card">
      <h3>Apache 2.0 license</h3>
      <p>No BSL restrictions. No usage fees. No control plane dependency. Fork it, sell it, run it however you want.</p>
    </div>
  </div>
</section>

<section class="section tradeoffs">
  <h2>What You Should Consider</h2>
  <p>KafScale is not a drop-in replacement for every Kafka workload. Here's when it fits and when it doesn't.</p>
  <div class="grid">
    <div class="card">
      <h3>KafScale is for you if</h3>
      <ul>
        <li>Latency of 200-500ms is acceptable</li>
        <li>You run ETL, logs, or async events</li>
        <li>You want minimal ops and no disk management</li>
        <li>Apache 2.0 licensing matters to you</li>
        <li>You prefer self-hosted over managed services</li>
      </ul>
    </div>
    <div class="card">
      <h3>KafScale is not for you if</h3>
      <ul>
        <li>You need sub-10ms latency</li>
        <li>You require exactly-once semantics (transactions)</li>
        <li>You rely on compacted topics</li>
        <li>You need native Iceberg integration</li>
        <li>You want a fully managed service</li>
      </ul>
    </div>
  </div>
  <div class="hero-actions">
    <a class="button secondary" href="/comparison/">See full comparison with alternatives</a>
  </div>
</section>

<section class="section">
  <h2>How KafScale works</h2>
  <div class="diagram">
    <svg viewBox="0 0 900 370" role="img" aria-label="KafScale architecture diagram">
      <defs>
        <marker id="arrow" markerWidth="10" markerHeight="10" refX="6" refY="3" orient="auto">
          <path d="M0,0 L0,6 L6,3 z" fill="var(--diagram-stroke)"></path>
        </marker>
      </defs>
      <rect x="20" y="20" width="860" height="220" rx="18" ry="18" fill="var(--diagram-fill)" stroke="var(--diagram-stroke)" stroke-width="2"></rect>
      <text x="40" y="55" font-size="16" font-weight="600" fill="var(--diagram-text)">Kubernetes cluster</text>
      <text x="70" y="82" font-size="12" fill="var(--diagram-label)">Stateless brokers (HPA, scale out/in)</text>

      <rect x="70" y="100" width="120" height="55" rx="10" fill="var(--diagram-accent)" stroke="var(--diagram-stroke)"></rect>
      <text x="100" y="133" font-size="13" fill="var(--diagram-text)">Broker 0</text>

      <rect x="210" y="100" width="120" height="55" rx="10" fill="var(--diagram-accent)" stroke="var(--diagram-stroke)"></rect>
      <text x="240" y="133" font-size="13" fill="var(--diagram-text)">Broker 1</text>

      <rect x="350" y="100" width="120" height="55" rx="10" fill="var(--diagram-accent)" stroke="var(--diagram-stroke)"></rect>
      <text x="380" y="133" font-size="13" fill="var(--diagram-text)">Broker N</text>

      <rect x="550" y="85" width="280" height="90" rx="12" fill="var(--diagram-fill)" stroke="var(--diagram-stroke)" stroke-width="1.5"></rect>
      <text x="565" y="108" font-size="13" font-weight="600" fill="var(--diagram-text)">etcd cluster</text>
      <text x="565" y="158" font-size="11" fill="var(--diagram-label)">Topics, offsets, group state</text>
      <circle cx="600" cy="130" r="12" fill="var(--diagram-accent)" stroke="var(--diagram-stroke)"/>
      <circle cx="660" cy="130" r="12" fill="var(--diagram-accent)" stroke="var(--diagram-stroke)"/>
      <circle cx="720" cy="130" r="12" fill="var(--diagram-accent)" stroke="var(--diagram-stroke)"/>

      <rect x="70" y="175" width="180" height="50" rx="10" fill="var(--diagram-fill)" stroke="var(--diagram-stroke)"></rect>
      <text x="115" y="205" font-size="12" font-weight="600" fill="var(--diagram-text)">Operator</text>

      <line x1="470" y1="127" x2="545" y2="127" stroke="var(--diagram-stroke)" stroke-width="1.5" marker-end="url(#arrow)"></line>
      <text x="480" y="118" font-size="10" fill="var(--diagram-label)">metadata</text>

      <line x1="280" y1="155" x2="280" y2="250" stroke="var(--diagram-stroke)" stroke-width="1.5" marker-end="url(#arrow)"></line>
      <text x="290" y="210" font-size="10" fill="var(--diagram-label)">segments</text>

      <rect x="20" y="260" width="860" height="90" rx="16" fill="var(--diagram-fill)" stroke="var(--diagram-stroke)" stroke-width="2"></rect>
      <text x="40" y="290" font-size="14" font-weight="600" fill="var(--diagram-text)">Amazon S3 (11 nines durability)</text>

      <rect x="60" y="300" width="300" height="40" rx="8" fill="var(--diagram-accent)" stroke="var(--diagram-stroke)"></rect>
      <text x="80" y="325" font-size="12" fill="var(--diagram-text)">Log segments + indexes</text>

      <rect x="540" y="300" width="300" height="40" rx="8" fill="var(--diagram-accent)" stroke="var(--diagram-stroke)"></rect>
      <text x="560" y="325" font-size="12" fill="var(--diagram-text)">etcd snapshots (backup)</text>

      <line x1="690" y1="175" x2="690" y2="295" stroke="var(--diagram-stroke)" stroke-width="1.5" marker-end="url(#arrow)"></line>
      <text x="700" y="240" font-size="10" fill="var(--diagram-label)">snapshots</text>
    </svg>
  </div>
  <div class="hero-actions">
    <a class="button secondary" href="/architecture/">See detailed architecture flows</a>
  </div>
</section>

<section class="section">
  <h2>Get running in minutes</h2>
  <div class="quickstart-flow">
    <div class="quickstart-card" markdown="1">

### 1. Install the operator

```bash
helm upgrade --install kafscale deploy/helm/kafscale \
  --namespace kafscale --create-namespace \
  --set operator.etcdEndpoints={} \
  --set operator.image.tag=latest
```

</div>
    <div class="quickstart-card" markdown="1">

### 2. Create a topic

```yaml
apiVersion: kafscale.io/v1alpha1
kind: KafscaleTopic
metadata:
  name: orders
spec:
  clusterRef: demo
  partitions: 3
```

</div>
    <div class="quickstart-card" markdown="1">

### 3. Produce and consume

```bash
kafka-console-producer \
  --bootstrap-server 127.0.0.1:9092 \
  --topic orders
```

</div>
  </div>
  <div class="hero-actions">
    <a class="button secondary" href="/quickstart/">Full quickstart guide</a>
  </div>
</section>

<section class="section">
  <h2>Production-grade operations</h2>
  <div class="grid">
    <div class="card">
      <h3>Prometheus metrics</h3>
      <p>S3 health state, produce/fetch throughput, consumer lag, etcd snapshot age. Grafana dashboards included.</p>
    </div>
    <div class="card">
      <h3>Horizontal scaling</h3>
      <p>Add brokers instantly. No partition rebalancing. HPA scales on CPU or custom metrics.</p>
    </div>
    <div class="card">
      <h3>Automated backups</h3>
      <p>Operator snapshots etcd to S3 on a schedule. One-command restore.</p>
    </div>
    <div class="card">
      <h3>Health gating</h3>
      <p>Brokers track S3 availability. Degraded and unavailable states prevent data loss.</p>
    </div>
  </div>
  <div class="hero-actions">
    <a class="button secondary" href="/operations/">Operations guide</a>
  </div>
</section>

<section class="section">
  <h2>Documentation</h2>
  <div class="grid">
    <div class="card">
      <h3>Protocol compatibility</h3>
      <p>21 Kafka APIs supported. Produce, Fetch, Metadata, consumer groups, and more.</p>
      <a class="button secondary" href="/protocol/">View API docs</a>
    </div>
    <div class="card">
      <h3>Storage format</h3>
      <p>Segment layout, index structure, S3 key paths, and cache architecture.</p>
      <a class="button secondary" href="/storage-format/">Explore storage</a>
    </div>
    <div class="card">
      <h3>Security</h3>
      <p>TLS configuration, S3 IAM policies, and the roadmap for SASL and ACLs.</p>
      <a class="button secondary" href="/security/">Security guide</a>
    </div>
  </div>
</section>

<section class="section backers">
  <h2>Backed by</h2>
  <p>KafScale is developed and maintained with support from <a href="https://scalytics.io" target="_blank" rel="noreferrer">Scalytics, Inc.</a> and <a href="https://novatechflow.com" target="_blank" rel="noreferrer">NovaTechFlow</a>.</p>
  <p>Apache 2.0 licensed. No CLA required. <a href="https://github.com/novatechflow/kafscale/blob/main/CONTRIBUTING.md">Contributions welcome</a>.</p>
</section>