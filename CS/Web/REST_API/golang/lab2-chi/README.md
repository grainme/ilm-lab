## Lab 2: Learning Notes

> A summary of key concepts and techniques from the lab.

### 1. Database Migrations (`golang-migrate`)

*   **Concept:** Schema as code. Provides version control for your database, making changes repeatable (`up`) and reversible (`down`).
*   **CLI vs. Library:**
    *   **CLI:** For manual creation/application during development (`migrate create`, `migrate up`).
    *   **Library:** For programmatic execution on application startup to ensure the DB schema is always in sync with the code.
*   **Key Insight:** The `golang-migrate` library requires its own `database` and `source` drivers to be blank-imported (`_ "..."`) to function. The source path must be a URL (e.g., `file://db/migrations`) - weird syntax (check out `main.go`)

### 2. SQLC: Type-Safe Data Access

*   **Concept:** Generates type-safe Go code from raw SQL queries.
*   **Technique:** Write standard SQL with special `-- name:` comments. `sqlc` generates the corresponding Go functions.
*   **Key Insight:** Automatically uses parameterized queries, which is a defense against **SQL Injection** by separating the query's logic from user-supplied data. (it does escape SQL's special chars)

### 3. Database Performance

*   **N+1 Query Problem:** Performance anti-pattern where 1 query fetches a list, and then N subsequent queries run in a loop to get related data. Solved by using a `JOIN` to fetch all required data in a single, efficient query.
*   **Indexing:** Creates a special lookup table for a column (e.g., `reviews.movie_id`) to make `WHERE` clauses and `JOIN` operations dramatically faster.

### 6. Docker for Development

*   **Multi-Stage Builds:** A `Dockerfile` can have multiple `FROM` statements. This is used to create a small final image containing only the compiled binary and necessary assets, discarding the larger build environment.
*   **Cache Invalidation:** `docker compose up` uses a cached image by default. `docker compose up --build` is used to force a rebuild after changing source code or the `Dockerfile`.
