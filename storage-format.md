---
layout: doc
title: Storage Format
description: Segment file format, index layout, and S3 key structure for Kafscale storage.
---

# Storage Format

## S3 key layout

```
s3://{bucket}/{namespace}/{topic}/{partition}/segment-{base_offset}.kfs
s3://{bucket}/{namespace}/{topic}/{partition}/segment-{base_offset}.index
```

Example:

```
s3://kafscale-data/production/orders/0/segment-00000000000000000000.kfs
s3://kafscale-data/production/orders/0/segment-00000000000000000000.index
```

## Segment file binary format

```
┌────────────────────────────────────────────────────────────────┐
│ Segment Header (32 bytes)                                      │
├────────────────────────────────────────────────────────────────┤
│ Magic Number        │ 4 bytes  │ 0x4B414653 ("KAFS")           │
│ Version             │ 2 bytes  │ Format version (1)            │
│ Flags               │ 2 bytes  │ Compression, etc.             │
│ Base Offset         │ 8 bytes  │ First offset in segment       │
│ Message Count       │ 4 bytes  │ Number of messages            │
│ Created Timestamp   │ 8 bytes  │ Unix millis                   │
│ Reserved            │ 4 bytes  │ Future use                    │
├────────────────────────────────────────────────────────────────┤
│ Message Batch 1                                                │
├────────────────────────────────────────────────────────────────┤
│ Message Batch 2                                                │
├────────────────────────────────────────────────────────────────┤
│ ...                                                            │
├────────────────────────────────────────────────────────────────┤
│ Segment Footer (16 bytes)                                      │
├────────────────────────────────────────────────────────────────┤
│ CRC32               │ 4 bytes  │ Checksum of all batches       │
│ Last Offset         │ 8 bytes  │ Last offset in segment        │
│ Footer Magic        │ 4 bytes  │ 0x454E4421 ("END!")           │
└────────────────────────────────────────────────────────────────┘
```

## Message batch format

```
┌────────────────────────────────────────────────────────────────┐
│ Batch Header (49 bytes)                                        │
├────────────────────────────────────────────────────────────────┤
│ Base Offset         │ 8 bytes  │ First offset in batch         │
│ Batch Length        │ 4 bytes  │ Total bytes in batch          │
│ Partition Leader    │ 4 bytes  │ Epoch of leader               │
│ Magic               │ 1 byte   │ 2 (Kafka compat)              │
│ CRC32               │ 4 bytes  │ Checksum of batch             │
│ Attributes          │ 2 bytes  │ Compression, timestamp type   │
│ Last Offset Delta   │ 4 bytes  │ Offset of last msg - base     │
│ First Timestamp     │ 8 bytes  │ Timestamp of first message    │
│ Max Timestamp       │ 8 bytes  │ Max timestamp in batch        │
│ Producer ID         │ 8 bytes  │ -1 (no idempotence)           │
│ Producer Epoch      │ 2 bytes  │ -1                            │
│ Base Sequence       │ 4 bytes  │ -1                            │
│ Record Count        │ 4 bytes  │ Number of records in batch    │
├────────────────────────────────────────────────────────────────┤
│ Record 1                                                       │
│ Record 2                                                       │
│ ...                                                            │
└────────────────────────────────────────────────────────────────┘
```

## Index file format

```
┌────────────────────────────────────────────────────────────────┐
│ Index Header (16 bytes)                                        │
├────────────────────────────────────────────────────────────────┤
│ Magic               │ 4 bytes  │ 0x494458 ("IDX")              │
│ Version             │ 2 bytes  │ 1                             │
│ Entry Count         │ 4 bytes  │ Number of index entries       │
│ Interval            │ 4 bytes  │ Messages between entries      │
│ Reserved            │ 2 bytes  │ Future use                    │
├────────────────────────────────────────────────────────────────┤
│ Entry 1: Offset (8 bytes) + Position (4 bytes)                 │
│ Entry 2: Offset (8 bytes) + Position (4 bytes)                 │
│ ...                                                            │
└────────────────────────────────────────────────────────────────┘
```

## Retention and lifecycle policy setup

Use bucket lifecycle rules to expire older segments and indexes based on your retention policy. Align retention with Kafka topic settings so producers and consumers have consistent expectations.
