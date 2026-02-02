## Lab 3: Redis — Learning Notes

> Concise reference of the main caching concepts and architectural patterns explored in this lab.

### Goal
Understand how a caching layer improves API performance by reducing database load and latency, while learning practical strategies for invalidation and resiliency.

---

## Core Concepts

Caching stores frequently accessed data in a fast in-memory store (Redis) instead of querying a disk-based database, resulting in lower latency and better system throughput.

### Cache-Aside Pattern
The application manages the cache lifecycle.

**Flow:**
1. Check cache first.  
2. **Hit:** return data immediately.  
3. **Miss:** query database → store result in cache → return data.

The cache sits “aside” the database, used only when beneficial.

### Cache Invalidation — *Delete, Don’t Update*
When data changes, delete the cache key rather than updating it.

**Why:**
- Simpler and less error-prone.
- Prevents partially updated or inconsistent objects.
- The next request repopulates the cache with fresh data (“self-healing”).

### Write-Behind Pattern
Best suited for high-frequency, non-critical writes (e.g., view counters).

**Flow:** write to cache → cache persists to DB asynchronously.

**Trade-off:** very low latency with a small risk of data loss if the cache fails before persistence.

### TTL & Eviction
- **TTL (Time-To-Live):** ensures stale data eventually expires and helps control memory usage.
- **Eviction Policies:** define what happens when memory is full. (there are other policies)
  - `volatile-lru`: evicts least recently used keys with TTL (good for caching).
  - `noeviction`: rejects writes when memory is exhausted.

**TTL Return Values**
- `> 0` — seconds until expiration  
- `-1` — key exists without expiration  
- `-2` — key does not exist  

### Graceful Degradation
The cache is an optimization, not a dependency.

If Redis is unavailable:
- bypass the cache  
- fetch from the database  
- keep the application running (with slower performance)

---

## Write Strategies Comparison

| Strategy        | Behavior | Pros | Cons |
|----------------|--------|------|------|
| **Write-Through** | Write to cache and DB synchronously | Strong consistency | Higher latency |
| **Write-Behind** | Write to cache, persist later | Very low latency | Small risk of data loss |
