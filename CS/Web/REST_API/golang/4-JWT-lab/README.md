## Lab 4: JWT Authentication — Learning Notes

> Concise reference of authentication concepts, password security, and token-based auth patterns.

### Goal
Understand how to secure API routes using JWT tokens, implement proper password storage with Argon2id, and learn the distinction between authentication and authorization.

---

## Core Concepts

### Authentication vs Authorization

| Concept | Question | HTTP Status on Failure |
|---------|----------|------------------------|
| **Authentication** | "Who are you?" | 401 Unauthorized |
| **Authorization** | "Are you allowed to do this?" | 403 Forbidden |

**Request lifecycle:**
```
Request → Authentication middleware → Authorization middleware → Handler
              ↓                            ↓
         401 if unknown              403 if not permitted
```

---

## Sessions vs JWT

### Session-Based Auth (Traditional)
Server stores state. Client holds a session ID in a cookie.

**Flow:**
1. User sends credentials
2. Server validates → creates session in memory/Redis
3. Server sends `Set-Cookie: session_id=abc123`
4. Browser sends cookie on every request
5. Server looks up session → identifies user

**Scaling problem:** Multiple servers behind a load balancer don't share memory. Solutions: sticky sessions, shared Redis storage, or session replication — all add complexity.

### JWT (JSON Web Token)
Client holds the state. Server is stateless.

**Structure:** `header.payload.signature` (each part is base64-encoded)

| Part | Contains | Example |
|------|----------|---------|
| Header | Algorithm, token type | `{"alg": "HS256", "typ": "JWT"}` |
| Payload | Claims (user data) | `{"user_id": 42, "role": "admin", "exp": 1735689600}` |
| Signature | HMAC of header + payload | Ensures token wasn't tampered with |

**Validation flow:**
1. Extract token from `Authorization: Bearer <token>` header
2. Split into header.payload.signature
3. Recompute signature using server's secret key
4. Compare signatures — if mismatch, reject
5. Check `exp` claim — if expired, reject
6. Extract claims → attach to request context

**Key insight:** Any server with the secret key can validate any token. No shared storage needed.

---

## Secret Key

The password only your server knows. Used to sign and verify tokens.

- **Where it lives:** Environment variable (`JWT_SECRET`). Never in code, never in git.
- **If leaked:** Attacker can forge tokens as any user. Rotate immediately.

---

## What NOT to Put in JWT Payload

The payload is **visible to anyone** — base64 is encoding, not encryption.

```bash
echo "eyJ1c2VyX2lkIjo0Mn0=" | base64 -d
# Output: {"user_id":42}
```

**Never include:**
- Passwords (hashed or plain)
- Credit card numbers, SSNs, PII
- API keys or secrets

**Safe to include:**
- `user_id`, `role`, `exp`, `iat`
- Email (debatable, but common)

**Principle:** Minimum data needed to identify user and permissions. Everything else, look up from DB.

---

## Password Hashing

### Encryption vs Hashing

| Encryption | Hashing |
|------------|---------|
| Two-way (can decrypt) | One-way (cannot reverse) |
| Needs a key | No key, just math |
| Used for: secrets you need back | Used for: passwords |

### The Rainbow Table Attack

If `hash("password123")` always produces the same output, attackers can:
1. Pre-compute hashes for millions of common passwords
2. Steal your database
3. Look up hashes → instant password recovery

### Salt: The Defense

A **salt** is random data added to the password before hashing.

```
hash("password123" + "random_salt_abc") → unique hash
hash("password123" + "random_salt_xyz") → different hash
```

- Each user gets a unique random salt
- Same password → different hashes for different users
- Rainbow tables become useless
- Salt is stored with the hash (not secret)

### Why Argon2id

| Algorithm | Speed | Attacker's Advantage |
|-----------|-------|---------------------|
| MD5, SHA-256 | Very fast | Billions of attempts/second on GPU |
| bcrypt | Slow | Thousands of attempts/second |
| Argon2id | Slow + memory-hard | ~1000 attempts/second, GPUs ineffective |

**Argon2id parameters:**
- `m` — memory in KB (e.g., 65536 = 64MB per hash)
- `t` — time/iterations
- `p` — parallelism (threads)

**Output format:** `$argon2id$v=19$m=65536,t=1,p=4$SALT$HASH`

**Principle:** You can't prevent database leaks. But you can make cracking so expensive that it's not worth the attacker's time.
