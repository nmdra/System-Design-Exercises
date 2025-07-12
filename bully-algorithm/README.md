# Bully Algorithm (Leader Election) with Docker & Redis

This project demonstrates a distributed **leader election** system using the **Bully Algorithm**, implemented in Go, running each node as a separate **Docker container**, and using **Redis Pub/Sub** for communication.

### 1. Build and Run

```bash
docker compose up --build
````

You will see logs from each node showing elections and leader announcements.

### 2. Change Node Count

```yaml
  node4:
    build: .
    environment:
      - NODE_ID=4
      - PEERS=1,2,3,4
      - REDIS_ADDR=redis:6379
```

> [!IMPORTANT]
> * Only the highest-ID node that is alive becomes the leader.
> * Nodes publish/subscribe on the `election` Redis channel.
> * Extendable for log replication or failure detection.
