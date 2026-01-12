---
layout: doc
title: Architecture
description: Architecture and data flows for KafScale brokers, metadata, and S3 segment storage.
permalink: /architecture/
---

# Architecture

KafScale brokers are stateless pods on Kubernetes. Metadata lives in etcd, while immutable log segments live in S3. Clients speak the Kafka protocol to a proxy that abstracts broker topology. Brokers flush segments to S3 and serve reads with caching.

---

## Platform overview

<div class="diagram">
  <div class="diagram-label">Architecture overview</div>
  <svg viewBox="0 0 850 420" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="KafScale architecture overview">
    <defs>
      <marker id="ah" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto"><path d="M0,0 L10,5 L0,10 z" fill="var(--diagram-stroke)"/></marker>
      <marker id="ao" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto"><path d="M0,0 L10,5 L0,10 z" fill="#ffb347"/></marker>
      <marker id="ag" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto"><path d="M0,0 L10,5 L0,10 z" fill="#4ab7f1"/></marker>
    </defs>

    <!-- Clients -->
    <rect x="30" y="30" width="150" height="55" rx="10" fill="var(--diagram-fill)" stroke="var(--diagram-stroke)" stroke-width="1.5"/>
    <text x="105" y="53" font-size="12" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Kafka Clients</text>
    <text x="105" y="70" font-size="10" fill="var(--diagram-label)" text-anchor="middle">producers &amp; consumers</text>

    <!-- Proxy -->
    <rect x="30" y="115" width="150" height="55" rx="10" fill="rgba(255, 179, 71, 0.15)" stroke="#ffb347" stroke-width="2"/>
    <text x="105" y="138" font-size="12" font-weight="600" fill="#ffb347" text-anchor="middle">Proxy</text>
    <text x="105" y="155" font-size="10" fill="var(--diagram-label)" text-anchor="middle">rewrites metadata</text>

    <!-- Clients to Proxy arrow -->
    <path d="M105,85 L105,110" stroke="var(--diagram-stroke)" stroke-width="2" fill="none" marker-end="url(#ah)"/>
    <text x="115" y="100" font-size="9" fill="var(--diagram-label)">single IP</text>

    <!-- Proxy to K8s arrow -->
    <path d="M180,142 L230,142" stroke="var(--diagram-stroke)" stroke-width="2" fill="none" marker-end="url(#ah)"/>

    <!-- K8s boundary -->
    <rect x="240" y="95" width="440" height="240" rx="14" fill="var(--diagram-accent)" stroke="#326ce5" stroke-width="2" stroke-dasharray="8,4"/>
    <text x="265" y="120" font-size="11" font-weight="600" fill="var(--diagram-label)" letter-spacing="0.5px">KUBERNETES CLUSTER</text>

    <!-- Brokers -->
    <rect x="270" y="140" width="100" height="60" rx="10" fill="rgba(106, 167, 255, 0.2)" stroke="#6aa7ff" stroke-width="1.5"/>
    <text x="320" y="165" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Broker 0</text>
    <text x="320" y="182" font-size="9" fill="var(--diagram-label)" text-anchor="middle">stateless</text>

    <rect x="385" y="140" width="100" height="60" rx="10" fill="rgba(106, 167, 255, 0.2)" stroke="#6aa7ff" stroke-width="1.5"/>
    <text x="435" y="165" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Broker 1</text>
    <text x="435" y="182" font-size="9" fill="var(--diagram-label)" text-anchor="middle">stateless</text>

    <rect x="500" y="140" width="100" height="60" rx="10" fill="rgba(106, 167, 255, 0.2)" stroke="#6aa7ff" stroke-width="1.5"/>
    <text x="550" y="165" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Broker N</text>
    <text x="550" y="182" font-size="9" fill="var(--diagram-label)" text-anchor="middle">stateless</text>

    <text x="620" y="165" font-size="9" fill="#6aa7ff">← HPA</text>

    <!-- etcd -->
    <rect x="310" y="250" width="200" height="65" rx="10" fill="rgba(74, 183, 241, 0.15)" stroke="#4ab7f1" stroke-width="1.5"/>
    <text x="410" y="278" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">etcd (3 nodes)</text>
    <text x="410" y="295" font-size="9" fill="var(--diagram-label)" text-anchor="middle">topics, offsets, group assignments</text>

    <!-- S3 Data -->
    <rect x="700" y="120" width="130" height="100" rx="12" fill="rgba(255, 179, 71, 0.12)" stroke="#ffb347" stroke-width="2"/>
    <circle cx="765" cy="160" r="24" fill="rgba(255, 179, 71, 0.3)" stroke="#ffb347" stroke-width="1.5"/>
    <text x="765" y="165" font-size="12" font-weight="700" fill="#ffb347" text-anchor="middle">S3</text>
    <text x="765" y="200" font-size="10" font-weight="600" fill="#ffb347" text-anchor="middle">Data Bucket</text>

    <!-- S3 Backup -->
    <rect x="700" y="250" width="130" height="65" rx="12" fill="rgba(255, 179, 71, 0.12)" stroke="#ffb347" stroke-width="2"/>
    <text x="765" y="278" font-size="10" font-weight="600" fill="#ffb347" text-anchor="middle">S3 Backup</text>
    <text x="765" y="295" font-size="9" fill="var(--diagram-label)" text-anchor="middle">etcd snapshots</text>

    <!-- Brokers to S3 -->
    <path d="M600,170 L695,160" stroke="#ffb347" stroke-width="2" fill="none" marker-end="url(#ao)"/>
    <text x="630" y="152" font-size="9" fill="#ffb347">flush segments</text>

    <path d="M695,178 L600,185" stroke="#6aa7ff" stroke-width="2" fill="none" marker-end="url(#ah)"/>
    <text x="630" y="198" font-size="9" fill="#6aa7ff">fetch + cache</text>

    <!-- Brokers to etcd -->
    <path d="M410,200 L410,245" stroke="#4ab7f1" stroke-width="1.5" fill="none" marker-end="url(#ag)"/>
    <text x="420" y="225" font-size="9" fill="var(--diagram-label)">metadata</text>

    <!-- etcd to S3 backup -->
    <path d="M510,282 L695,282" stroke="#ffb347" stroke-width="1.5" fill="none" marker-end="url(#ao)"/>
    <text x="590" y="272" font-size="9" fill="var(--diagram-label)">snapshots</text>

    <!-- Tagline -->
    <text x="425" y="390" font-size="10" fill="var(--diagram-label)" text-anchor="middle" font-style="italic">
      One endpoint for clients. S3 is the source of truth. Brokers are stateless.
    </text>
  </svg>
</div>

---

## How the proxy works

The Kafka protocol requires clients to discover broker topology. When a client connects, the broker returns a list of all brokers and their partition assignments. Clients then connect directly to each broker they need.

This creates a problem for ephemeral infrastructure. Every broker restart breaks client connections. Scaling events require clients to rediscover the cluster.

KafScale's proxy solves this by intercepting two types of responses:

| Request | What the proxy does |
|---------|---------------------|
| **Metadata** | Returns the proxy's own address instead of individual broker addresses |
| **FindCoordinator** | Returns the proxy's address for consumer group coordination |

Clients believe they are talking to a single broker. The proxy routes requests to the actual brokers internally.

This enables:

- **Infinite horizontal scaling**: Add brokers without client awareness
- **Zero-downtime deployments**: Rotate broker pods behind the proxy
- **Standard networking**: One LoadBalancer, one DNS name, standard TLS termination

For configuration details, see [Operations: External Broker Access](/operations/#external-broker-access).

---

## Decoupled processing (addons)

KafScale keeps brokers focused on Kafka protocol and storage. Add-on processors handle downstream tasks by reading completed segments directly from S3, bypassing brokers entirely. Processors are stateless: offsets and leases live in etcd, input lives in S3, output goes to external catalogs.

<div class="diagram">
  <div class="diagram-label">Data processor architecture</div>
  <svg viewBox="0 0 800 320" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="KafScale processor addon architecture">
    <defs>
      <marker id="ap1" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto">
        <path d="M0,0 L10,5 L0,10 z" fill="var(--diagram-stroke)"/>
      </marker>
      <marker id="ap2" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto">
        <path d="M0,0 L10,5 L0,10 z" fill="#ffb347"/>
      </marker>
      <marker id="ap3" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto">
        <path d="M0,0 L10,5 L0,10 z" fill="#34d399"/>
      </marker>
      <marker id="ap4" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto">
        <path d="M0,0 L10,5 L0,10 z" fill="#4ab7f1"/>
      </marker>
      <marker id="ap5" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto">
        <path d="M0,0 L10,5 L0,10 z" fill="#a78bfa"/>
      </marker>
    </defs>

    <!-- Row 1: KafScale layer (source) -->
    <text x="30" y="28" font-size="10" font-weight="600" fill="var(--diagram-label)" letter-spacing="0.5px">KAFSCALE</text>
    
    <rect x="30" y="40" width="130" height="60" rx="10" fill="rgba(106, 167, 255, 0.2)" stroke="#6aa7ff" stroke-width="1.5"/>
    <text x="95" y="65" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Brokers</text>
    <text x="95" y="82" font-size="9" fill="var(--diagram-label)" text-anchor="middle">Kafka protocol</text>

    <rect x="190" y="40" width="130" height="60" rx="10" fill="rgba(255, 179, 71, 0.12)" stroke="#ffb347" stroke-width="2"/>
    <text x="255" y="65" font-size="11" font-weight="700" fill="#ffb347" text-anchor="middle">S3</text>
    <text x="255" y="82" font-size="9" fill="var(--diagram-label)" text-anchor="middle">.kfs segments</text>

    <rect x="350" y="40" width="130" height="60" rx="10" fill="rgba(74, 183, 241, 0.15)" stroke="#4ab7f1" stroke-width="1.5"/>
    <text x="415" y="65" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">etcd</text>
    <text x="415" y="82" font-size="9" fill="var(--diagram-label)" text-anchor="middle">metadata, offsets</text>

    <!-- Broker to S3 arrow -->
    <path d="M160,70 L185,70" stroke="#ffb347" stroke-width="1.5" fill="none" marker-end="url(#ap2)"/>

    <!-- Row 2: Processor layer -->
    <text x="30" y="138" font-size="10" font-weight="600" fill="var(--diagram-label)" letter-spacing="0.5px">PROCESSOR</text>

    <rect x="130" y="150" width="260" height="70" rx="12" fill="rgba(52, 211, 153, 0.12)" stroke="#34d399" stroke-width="2"/>
    <text x="260" y="175" font-size="12" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Processor</text>
    <text x="260" y="195" font-size="9" fill="var(--diagram-label)" text-anchor="middle">stateless pods, topic-scoped leases, HPA</text>

    <!-- S3 to Processor arrow -->
    <path d="M255,100 L255,145" stroke="#ffb347" stroke-width="2" fill="none" marker-end="url(#ap2)"/>
    <text x="268" y="125" font-size="9" fill="#ffb347">read segments</text>

    <!-- etcd to Processor arrow -->
    <path d="M415,100 L415,130 L360,130 L360,145" stroke="#4ab7f1" stroke-width="1.5" fill="none" marker-end="url(#ap4)"/>
    <text x="420" y="125" font-size="9" fill="#4ab7f1">offsets, leases</text>

    <!-- Row 3: Output layer -->
    <text x="30" y="258" font-size="10" font-weight="600" fill="var(--diagram-label)" letter-spacing="0.5px">OUTPUT</text>

    <rect x="80" y="270" width="150" height="40" rx="10" fill="rgba(167, 139, 250, 0.15)" stroke="#a78bfa" stroke-width="1.5"/>
    <text x="155" y="295" font-size="10" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Metadata Catalog</text>

    <rect x="260" y="270" width="130" height="40" rx="10" fill="rgba(255, 179, 71, 0.12)" stroke="#ffb347" stroke-width="1.5"/>
    <text x="325" y="295" font-size="10" font-weight="600" fill="#ffb347" text-anchor="middle">S3 Warehouse</text>

    <!-- Processor to Catalog arrow -->
    <path d="M200,220 L155,265" stroke="#a78bfa" stroke-width="2" fill="none" marker-end="url(#ap5)"/>
    <text x="145" y="245" font-size="9" fill="#a78bfa">metadata</text>

    <!-- Processor to Warehouse arrow -->
    <path d="M310,220 L325,265" stroke="#ffb347" stroke-width="2" fill="none" marker-end="url(#ap2)"/>
    <text x="340" y="245" font-size="9" fill="#ffb347">parquet</text>

    <!-- Consumers on right side -->
    <text x="520" y="258" font-size="10" font-weight="600" fill="var(--diagram-label)" letter-spacing="0.5px">CONSUMERS</text>

    <rect x="520" y="270" width="130" height="40" rx="10" fill="var(--diagram-fill)" stroke="var(--diagram-stroke)" stroke-width="1.5"/>
    <text x="585" y="290" font-size="9" fill="var(--diagram-text)" text-anchor="middle">Analytics, AI agents</text>
    <text x="585" y="302" font-size="9" fill="var(--diagram-label)" text-anchor="middle">query engines</text>

    <!-- Catalog/Warehouse to Consumers -->
    <path d="M390,290 L515,290" stroke="var(--diagram-stroke)" stroke-width="1.5" stroke-dasharray="4,2" fill="none" marker-end="url(#ap1)"/>
    <text x="452" y="282" font-size="9" fill="var(--diagram-label)">query</text>

    <!-- Caption -->
    <text x="400" y="30" font-size="10" fill="var(--diagram-label)" text-anchor="middle" font-style="italic">
      Processors bypass brokers entirely. State lives in etcd. Data lives in S3.
    </text>
  </svg>
</div>

The processor reads .kfs segments from S3, tracks progress in etcd, and writes Parquet files to an Iceberg warehouse. Any Iceberg-compatible catalog can serve the tables to downstream consumers.

For deployment and configuration, see the [Iceberg Processor](/processors/iceberg/) docs.

---

## Key design decisions

| Decision | Rationale |
|----------|-----------|
| **Proxy for topology abstraction** | Clients see one endpoint. Brokers scale without client awareness. |
| **S3 as source of truth** | 11 nines durability, unlimited capacity, ~$0.023/GB/month |
| **Stateless brokers** | Any pod serves any partition. HPA scales 0→N instantly. |
| **etcd for metadata** | Leverages existing K8s patterns. Strong consistency. |
| **~500ms latency** | Acceptable trade-off for ETL, logs, async events |
| **No transactions** | Simplifies architecture. Covers 80% of Kafka use cases. |
| **4MB segment size** | Balances S3 PUT costs (~$0.005/1000) vs flush latency |

---

## Produce flow

<div class="diagram">
  <div class="diagram-label">Write path</div>
  <svg viewBox="0 0 750 130" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="KafScale produce flow">
    <defs>
      <marker id="ab" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto"><path d="M0,0 L10,5 L0,10 z" fill="#6aa7ff"/></marker>
      <marker id="ao2" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto"><path d="M0,0 L10,5 L0,10 z" fill="#ffb347"/></marker>
    </defs>

    <!-- Producer -->
    <rect x="20" y="30" width="120" height="70" rx="10" fill="var(--diagram-fill)" stroke="var(--diagram-stroke)" stroke-width="1.5"/>
    <text x="80" y="58" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Producer</text>
    <text x="80" y="75" font-size="9" fill="var(--diagram-label)" text-anchor="middle">Kafka client</text>

    <!-- Broker -->
    <rect x="200" y="30" width="140" height="70" rx="10" fill="rgba(106, 167, 255, 0.2)" stroke="#6aa7ff" stroke-width="1.5"/>
    <text x="270" y="55" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Broker</text>
    <text x="270" y="72" font-size="9" fill="var(--diagram-label)" text-anchor="middle">validate, batch</text>
    <text x="270" y="86" font-size="9" fill="var(--diagram-label)" text-anchor="middle">assign offsets</text>

    <!-- Buffer -->
    <rect x="400" y="30" width="130" height="70" rx="10" fill="rgba(52, 211, 153, 0.15)" stroke="#34d399" stroke-width="1.5"/>
    <text x="465" y="58" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Buffer</text>
    <text x="465" y="75" font-size="9" fill="var(--diagram-label)" text-anchor="middle">in-memory batches</text>

    <!-- S3 -->
    <rect x="590" y="30" width="130" height="70" rx="12" fill="rgba(255, 179, 71, 0.12)" stroke="#ffb347" stroke-width="2"/>
    <circle cx="655" cy="58" r="18" fill="rgba(255, 179, 71, 0.3)" stroke="#ffb347" stroke-width="1"/>
    <text x="655" y="63" font-size="11" font-weight="700" fill="#ffb347" text-anchor="middle">S3</text>
    <text x="655" y="85" font-size="9" fill="var(--diagram-label)" text-anchor="middle">sealed segment</text>

    <!-- Arrows -->
    <path d="M140,65 L195,65" stroke="#6aa7ff" stroke-width="2" fill="none" marker-end="url(#ab)"/>
    <text x="165" y="55" font-size="10" font-weight="600" fill="#6aa7ff">1</text>

    <path d="M340,65 L395,65" stroke="#6aa7ff" stroke-width="2" fill="none" marker-end="url(#ab)"/>
    <text x="365" y="55" font-size="10" font-weight="600" fill="#6aa7ff">2</text>

    <path d="M530,65 L585,65" stroke="#ffb347" stroke-width="2" fill="none" marker-end="url(#ao2)"/>
    <text x="555" y="55" font-size="10" font-weight="600" fill="#ffb347">3</text>

    <!-- Labels -->
    <text x="167" y="118" font-size="9" fill="var(--diagram-label)" text-anchor="middle">produce</text>
    <text x="367" y="118" font-size="9" fill="var(--diagram-label)" text-anchor="middle">batch</text>
    <text x="557" y="118" font-size="9" fill="var(--diagram-label)" text-anchor="middle">flush</text>
  </svg>
</div>

1. **Produce**: Client sends records to any broker via Kafka protocol
2. **Batch**: Broker validates, batches records, assigns offsets
3. **Flush**: When buffer reaches 4MB or 500ms, segment is sealed and uploaded to S3

Data is not acknowledged until S3 upload completes. This guarantees 11 nines durability on ACK.

---

## Fetch flow

<div class="diagram">
  <div class="diagram-label">Read path</div>
  <svg viewBox="0 0 750 160" xmlns="http://www.w3.org/2000/svg" role="img" aria-label="KafScale fetch flow">
    <defs>
      <marker id="ab3" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto"><path d="M0,0 L10,5 L0,10 z" fill="#6aa7ff"/></marker>
      <marker id="ao3" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="5" markerHeight="5" orient="auto"><path d="M0,0 L10,5 L0,10 z" fill="#ffb347"/></marker>
    </defs>

    <!-- Consumer -->
    <rect x="20" y="35" width="120" height="70" rx="10" fill="var(--diagram-fill)" stroke="var(--diagram-stroke)" stroke-width="1.5"/>
    <text x="80" y="63" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Consumer</text>
    <text x="80" y="80" font-size="9" fill="var(--diagram-label)" text-anchor="middle">Kafka client</text>

    <!-- Broker -->
    <rect x="200" y="35" width="140" height="70" rx="10" fill="rgba(106, 167, 255, 0.2)" stroke="#6aa7ff" stroke-width="1.5"/>
    <text x="270" y="60" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">Broker</text>
    <text x="270" y="77" font-size="9" fill="var(--diagram-label)" text-anchor="middle">locate segment</text>
    <text x="270" y="91" font-size="9" fill="var(--diagram-label)" text-anchor="middle">check cache</text>

    <!-- Cache -->
    <rect x="400" y="35" width="130" height="70" rx="10" fill="rgba(56, 189, 248, 0.15)" stroke="#38bdf8" stroke-width="1.5"/>
    <text x="465" y="63" font-size="11" font-weight="600" fill="var(--diagram-text)" text-anchor="middle">LRU Cache</text>
    <text x="465" y="80" font-size="9" fill="var(--diagram-label)" text-anchor="middle">hit → fast path</text>

    <!-- S3 -->
    <rect x="590" y="35" width="130" height="70" rx="12" fill="rgba(255, 179, 71, 0.12)" stroke="#ffb347" stroke-width="2"/>
    <circle cx="655" cy="63" r="18" fill="rgba(255, 179, 71, 0.3)" stroke="#ffb347" stroke-width="1"/>
    <text x="655" y="68" font-size="11" font-weight="700" fill="#ffb347" text-anchor="middle">S3</text>
    <text x="655" y="90" font-size="9" fill="var(--diagram-label)" text-anchor="middle">miss → fetch</text>

    <!-- Forward arrows -->
    <path d="M140,70 L195,70" stroke="#6aa7ff" stroke-width="2" fill="none" marker-end="url(#ab3)"/>
    <text x="165" y="60" font-size="10" font-weight="600" fill="#6aa7ff">1</text>

    <path d="M340,70 L395,70" stroke="#6aa7ff" stroke-width="2" fill="none" marker-end="url(#ab3)"/>
    <text x="365" y="60" font-size="10" font-weight="600" fill="#6aa7ff">2</text>

    <path d="M530,70 L585,70" stroke="#ffb347" stroke-width="2" fill="none" marker-end="url(#ao3)"/>
    <text x="555" y="60" font-size="10" font-weight="600" fill="#ffb347">3</text>

    <!-- Return arrows -->
    <path d="M585,90 Q465,145 340,90" stroke="#ffb347" stroke-width="1.5" stroke-dasharray="4,2" fill="none" marker-end="url(#ao3)"/>
    <text x="465" y="140" font-size="10" font-weight="600" fill="#ffb347">4</text>

    <path d="M195,90 Q130,130 80,105" stroke="#6aa7ff" stroke-width="1.5" stroke-dasharray="4,2" fill="none" marker-end="url(#ab3)"/>
    <text x="125" y="135" font-size="10" font-weight="600" fill="#6aa7ff">5</text>

    <!-- Labels -->
    <text x="167" y="25" font-size="9" fill="var(--diagram-label)" text-anchor="middle">fetch</text>
    <text x="367" y="25" font-size="9" fill="var(--diagram-label)" text-anchor="middle">cache?</text>
    <text x="557" y="25" font-size="9" fill="var(--diagram-label)" text-anchor="middle">miss</text>
  </svg>
</div>

1. **Fetch**: Consumer requests data from broker
2. **Cache check**: Broker looks up segment in LRU cache
3. **S3 fetch**: On cache miss, broker fetches from S3
4. **Populate**: Fetched segment is cached for future requests
5. **Return**: Data returned to consumer

---

## Component responsibilities

| Component | Responsibilities |
|-----------|------------------|
| **Proxy** | Rewrites Metadata/FindCoordinator responses, routes requests to brokers, enables topology abstraction |
| **Broker** | Kafka protocol, batching, offset assignment, S3 read/write, caching |
| **etcd** | Topic metadata, consumer offsets, group assignments, leader election |
| **S3** | Durable segment storage, source of truth, lifecycle-based retention |
| **Operator** | CRD reconciliation, etcd snapshots, broker lifecycle management |

---

## Segment format summary

Segments are self-contained files with header, Kafka-compatible record batches, and footer.

| Field | Size | Description |
|-------|------|-------------|
| Magic | 4 bytes | `0x4B414653` ("KAFS") |
| Version | 2 bytes | Format version (1) |
| Flags | 2 bytes | Compression codec |
| Base Offset | 8 bytes | First offset in segment |
| Message Count | 4 bytes | Number of messages |
| Timestamp | 8 bytes | Created (Unix ms) |
| Batches | variable | Kafka RecordBatch format |
| CRC32 | 4 bytes | Checksum |
| Footer Magic | 4 bytes | `0x454E4421` ("END!") |

See [Storage Format](/storage-format/) for complete details on segment structure, indexes, and S3 key layout.

---

## Next steps

- [Operations](/operations/) for proxy configuration and S3 health states
- [Storage Format](/storage-format/) for detailed segment and index layouts
- [Rationale](/rationale/) for why we made these design choices