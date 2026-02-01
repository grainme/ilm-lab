## Lab 3: Redis Learning Notes

> A summary of key concepts and architectural patterns from our Redis mentoring session.

### 1. The Cache-Aside Pattern

This is the fundamental caching logic we designed. The application, not the cache, is responsible for managing data.

-   **Logic Flow**:
    1.  The application first asks the cache for the data.
    2.  **Cache Hit**: If the data exists, it's returned immediately.
    3.  **Cache Miss**: If the data does not exist, the application queries the database, stores the result in the cache for next time, and then returns the data.
-   **Why the name "cache-aside"?**: The cache sits "to the side" of the database. The application talks to the cache, and if that fails, it goes around the cache to talk directly to the database.

### 2. Cache Invalidation: Delete, Don't Update

When data changes in the database (due to an update or delete), we must handle the stale data in the cache.

-   **The Strategy**: Instead of *updating* the data in the cache, the simpler and more robust approach is to just **`DEL`** the key from the cache.
-   **Why it's better**:
    -   **Simplicity**: Your `UPDATE` and `DELETE` logic becomes identical: just delete the key.
    -   **Self-Healing**: The next time the data is requested, it's a cache miss. The cache-aside logic then automatically fetches the fresh data from the database and repopulates the cache. This prevents complex bugs where you might only partially update the cached object.

### 3. The Critical Role of `EXPIRE`

Setting an expiration on your cache keys is crucial for two main reasons.

1.  **A Safety Net**: It guarantees that stale data will eventually be removed, even if your active `DEL`-on-update logic fails due to a bug or a crash.
2.  **Memory Management**: It acts as an automatic garbage collector, ensuring your cache doesn't grow indefinitely and run out of RAM.

### 4. Key Redis Internals You Asked About

-   **`TTL` Return Values**:
    -   `> 0`: The number of seconds until expiration.
    -   `-1`: The key exists but has no expiration set (it's persistent).
    -   `-2`: The key does not exist (it was deleted or has already expired).

-   **Eviction Policies**:
    -   The default policy is **`noeviction`**, which returns an error on writes when memory is full.
    -   You must manually configure policies like **`volatile-lru`** (a great choice for caching, as it only evicts keys that have an expiration set) or `allkeys-lru`.

### 5. Docker Networking for Development

-   **The Problem**: A Go app running on your host machine can't connect to `localhost:6379` to reach a service inside a Docker container, because they are in different network worlds.
-   **The Solutions**:
    1.  **Port Mapping**: Use `docker run -p 6379:6379 ...` to create a portal from your host's `localhost` into the container. This is for when your app runs *outside* Docker.
    2.  **Container Networking**: Run your app *inside* a Docker container using `docker-compose`. The containers can then talk to each other using their service names as hostnames (e.g., your Go app connects to `redis:6379`). This is the recommended approach for multi-service applications.