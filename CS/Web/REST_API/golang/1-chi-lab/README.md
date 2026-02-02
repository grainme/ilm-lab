Lab 1: Learning Summary

#### 1. Layered Architecture
-   **Separation of Concerns**: We have a three-layer architecture (`handler`, `service`, `repository`) which is a standard pattern for building Go applications.
    -   **Handlers** are responsible only for HTTP request/response logic.
    -   **Services** contain the core business logic and validation.
    -   **Repositories** are responsible only for data access from data source (DB, Cache, Files, In-Memory...)
-   **DI & Interfaces**: By defining a `MovieRepository` interface, the `service` layer depends on an abstraction, not a concrete implementation. This allows the data source to be swapped (e.g., from our in-memory store to the Postgres DB in Lab 2) without changing the service code.

#### 2. Data Race Safety
-   **`sync.RWMutex`**: We used a read/write mutex to protect shared data (`movies` slice) from concurrent access, using `RLock` for reads and `Lock` for writes.
-   **Defensive Copying**: The most critical lesson here. We learned that returning a pointer (or a slice of pointers) to internal, mutable state creates a **data race** vulnerability. The caller can modify the data without acquiring the lock.
    -   **Solution**: Always return a *copy* of the data. This ensures the repository maintains sole ownership and control over its internal state. This applies to both read (`Get...`) and write (`Add...`, `Update...`) methods that return data.

#### 3. Context
-   **`context.Context` Propagation**: We plumbed `context.Context` through all layers of the application. While not fully utilized in this lab, it's a foundational pattern in Go for handling request timeouts and cancellations, which will be vital in future labs involving network I/O (like database calls in Lab 2).
