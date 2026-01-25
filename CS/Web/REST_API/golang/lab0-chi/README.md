## Lab 0: Docker Basics - Containerize Your Development

**Goal**: Learn Docker for running Postgres/Redis locally  
**Spring analogy**: Docker = portable JVM, containers = lightweight VMs

**Theory Covered**:
- **Container vs VM**: Container shares OS kernel (lighter), VM has full OS (heavier)
- **Image vs Container**: Image = blueprint, Container = running instance
- **Docker Compose**: Run multiple containers together (Postgres + Redis)
- **Volumes**: Persist data after container restarts

### Project: Run Postgres in Docker (for Movie API)

### Tasks
1. **Install Docker Desktop**:
   - Download from docker.com
   - Verify: `docker --version`

2. **Understand key concepts**:
   - **Image**: Packaged software (like `postgres:15`)
   - **Container**: Running image (like running Spring Boot JAR)
   - **Volume**: Persistent storage (survives container deletion)
   - **Network**: Containers talk to each other

3. **Create docker-compose.yml in project root**:
   ```
   lab1-chi/
   ├── docker-compose.yml    # NEW
   ├── cmd/
   └── internal/
   ```

4. **Define services**:
   - Postgres (port 5432)
   - Redis (port 6379) - we'll use later
   - pgAdmin (web UI for Postgres, optional)

5. **Learn Docker commands**:
   - `docker-compose up -d` - start containers in background
   - `docker-compose down` - stop and remove containers
   - `docker-compose logs postgres` - view logs
   - `docker ps` - list running containers
   - `docker exec -it <container> psql -U postgres` - connect to Postgres

6. **Test connection**:
   - Start containers
   - Connect with `psql` or pgAdmin
   - Create test database: `CREATE DATABASE movies_dev;`

### Structure
```
lab1-chi/
├── docker-compose.yml       # NEW - defines Postgres + Redis
├── .env.example             # NEW - example environment variables
├── cmd/
└── internal/
```

### Learning Resources
- [Docker tutorial](https://docs.docker.com/get-started/)
- [Docker Compose overview](https://docs.docker.com/compose/)
- Focus on: `docker-compose up/down`, port mapping, volumes

### Success Criteria
- ✅ `docker-compose up -d` starts Postgres
- ✅ Can connect to Postgres at `localhost:5432`
- ✅ Data persists after `docker-compose down` and `up` (volume works)
- ✅ Can view logs with `docker-compose logs`

**Theory Test**: 
1. What happens to data if you run `docker-compose down -v`? (Answer: -v deletes volumes, data lost)
2. Why use Docker instead of installing Postgres directly?
