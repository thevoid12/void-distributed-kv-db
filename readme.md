# Void-KV: Distributed Key-Value Store

Void-KV is a distributed key-value store implemented on top of bolt db which includes static sharding with read replication across zones, consistent hashing for efficient data distribution, and automatic re-sharding. The project also includes comprehensive performance benchmarking for read and write operations.

## Features

- **Static Sharding with Read Replication**: Distributes data across multiple shards with dedicated read replicas in different zones to ensure high availability and reduce latency.
- **Consistent Hashing**: Determines data placement and shard assignments, enabling balanced data distribution and minimal data movement during re-sharding.
- **Automatic Re-sharding**: Adjusts shard allocation in response to load changes, supporting scalability and data rebalancing.
- **Performance Benchmarking**: Built-in benchmarking tools provide metrics for read and write performance across zones, helping identify optimizations and potential bottlenecks.

## Technology Stack

- **Programming Language**: Go
- **Networking**: gRPC- http gin(for inter-shard and replica communication)
- **Database**: boltdb
- **Load Balancing**: Consistent hashing ring
- **Benchmarking**: Custom scripts 

## Architecture

The architecture is built around a consistent hashing-based sharding mechanism, allowing each shard to manage a portion of the key space. Each shard includes:

1. **Primary Node**: Handles writes and maintains the authoritative state of the data.
2. **Replica Nodes**: Distributed across zones for read redundancy and low-latency data access.
3. **Consistent Hashing**: Ensures minimal reallocation of keys on shard or replica failure and enables efficient re-sharding.

## Demo video
[![void-kv-distributed db implementation](https://img.youtube.com/vi/IQgWHKAcaR0/hqdefault.jpg )](https://youtu.be/IQgWHKAcaR0?si=Txurosfka6GhI2sk)

## Implementation Details Explanation video
[![void-kv-distributed db explanation](https://img.youtube.com/vi/2mMF4qOHOgs/hqdefault.jpg )](https://youtu.be/2mMF4qOHOgs?si=5gH0Owoh8Q2Oqcsf)

## Getting Started

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/thevoid12/void-distributed-kv-db.git
   cd void-distributed-kv-db
   ```
2. setting up the dependency:
   ```bash
   go mod tidy
   ```
3. everything is encapsurated into the makefile: change db names, port and other customization from the make file
   ```bash
   make
   ```
  
   