# Kafscale User Guide

Kafscale is a Kafka-compatible, S3-backed message transport system. It keeps brokers stateless, stores data in S3, and relies on Kubernetes for scheduling and scaling. This guide summarizes how to interact with the platform once it is deployed.

## Concepts

- **Topics / Partitions**: match upstream Kafka semantics. All Kafka client libraries continue to work.
- **Brokers**: stateless pods accepting Kafka protocol traffic on port 9092 and metrics + gRPC control on 9093.
- **Metadata**: stored in etcd, encoded via protobufs (`kafscale.metadata.*`).
- **Storage**: message segments live in S3 buckets; brokers only keep in-memory caches.
- **Operator**: Kubernetes controller that provisions brokers, topics, and wiring based on CRDs.

## Getting Started

1. **Deploy the Operator + Brokers**  
   The Helm chart under `deploy/helm/kafscale` installs the CRDs, operator, and a broker StatefulSet. (Helm packaging wired up once the codebase is ready.)

2. **Create a Topic**  
   Apply a `KafscaleTopic` custom resource (see `config/samples/`). The operator writes the protobuf topic config into etcd; brokers pick it up automatically.

3. **Produce / Consume**  
   Point any Kafka client at the broker service:
   ```bash
   kafka-console-producer --bootstrap-server kafscale-broker:9092 --topic orders
   kafka-console-consumer --bootstrap-server kafscale-broker:9092 --topic orders --from-beginning
   ```

4. **Monitoring**  
   - Metrics via Prometheus on port 9093 (`/metrics`)
   - Structured JSON logs from brokers/operators
   - Control-plane queries via the gRPC service defined in `proto/control/broker.proto`

5. **Scaling / Maintenance**  
   The operator uses Kubernetes HPA and the BrokerControl gRPC API to safely drain partitions before restarts. Users can request manual drains or flushes by invoking those RPCs (CLI tooling TBD).

## Limits / Non-Goals

- No embedded stream processing features—pair Kafscale with Flink, Wayang, Spark, etc.
- Transactions, idempotent producers, and log compaction are out of scope for the MVP.

For deeper architectural details or development guidance, read `kscale-spec.md` and `docs/development.md`.
