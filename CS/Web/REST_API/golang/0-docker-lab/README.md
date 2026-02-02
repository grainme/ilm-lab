## Lab 0: Docker — Learning Notes

> Reference of foundational Docker concepts.

### Goal
Provide a reproducible local database environment without installing Postgres directly, avoiding version conflicts, unnecessary resource usage, and leftover system artifacts.

### Why Docker?
Docker offers **isolated, reproducible, and disposable** environments.

**Problems it solves:**
- Version mismatches across machines and production.
- Background services consuming CPU/RAM.
- Difficult cleanup after uninstalling software.

---

## Core Concepts

### Image vs. Container
- **Image:** Immutable blueprint containing the application and its dependencies.  
- **Container:** A running instance of an image.

### Docker Hub & Version Pinning
Use official images and pin exact versions instead of `latest` to ensure consistent builds and predictable behavior.

### Manual Container Management
Running containers manually helps understand ports, environment variables, and naming before automating workflows.

### Ephemeral Filesystem & Volumes
Containers lose their data when removed.

**Solution:** use volumes to persist data beyond the container lifecycle.

### Debugging Containers
- Check logs first when a container fails.
- Verify assumptions (paths, ports, configs).
- Small typos can silently prevent startup.

### Docker Compose
Defines multi-service applications declaratively and ensures services can communicate reliably. Ideal for setups like an API connected to a database.

---

## Multi-Stage Builds — Mental Model

### Only the Final Stage Ships
Docker may build multiple stages, but only the **last stage** becomes the production image.

**Why it matters:**
- Produces smaller images.
- Removes build tools from the final artifact.
- Improves security and startup time.

### Single vs. Multi-Stage
- **Single-stage:** build tools remain inside the final image.
- **Multi-stage:** compile in one stage, ship only the minimal runtime.

---

## Networking (Local Development)

**Challenge:** Host applications cannot directly access container services via `localhost`.

**Approaches:**
- **Port Mapping:** exposes container ports to the host.
- **Container Networking (Recommended):** services communicate using their names within a shared network.

---

## Next Concepts to Explore (other labs maybe TODO)
- **Healthchecks:** ensure services are truly ready before dependencies start.
- **.dockerignore:** exclude unnecessary files to keep images small and secure.
- **Bind Mounts:** enable live code updates during development.
- **Custom Networks:** improve service isolation.
- **Secrets Management:** inject sensitive data securely instead of storing it in configuration files.

---

# Go Static Binary — Learning Notes

> Understanding static vs. dynamic binaries and how Go handles dependencies.

## Definition
A **static binary** contains all required dependencies within a single executable, removing the need for external runtime libraries.

## Go’s Default Behavior
- **Pure Go code:** typically produces static binaries.
- **cgo usage:** introduces dynamic linking to system libraries.

### When cgo Is Triggered
- Certain networking behavior on Linux.
- Database drivers with C bindings.
- Any package that imports C code.
