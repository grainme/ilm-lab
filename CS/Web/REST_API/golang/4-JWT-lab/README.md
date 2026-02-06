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

---

## Middleware Implementation

### What Middleware Is

A function that wraps a handler. Runs before (and optionally after) your handler executes.

```
Request → [Middleware 1] → [Middleware 2] → [Handler] → Response
              ↓                  ↓
          (logging)         (auth check)
```

### The Signature

```go
func(next http.Handler) http.Handler
```

You receive the next handler, return a new handler that does something before/after calling `next`.

### Skeleton Pattern

```go
func Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. BEFORE: Extract token, validate
        // 2. DECISION: Call next or return early with error
        // 3. Pass modified request with context
        next.ServeHTTP(w, r.WithContext(newCtx))
    })
}
```

---

## Context for Request-Scoped Data

### Two Uses of Context

| Use | Purpose | Example |
|-----|---------|---------|
| Cancellation | "Should I stop?" | Timeouts, client disconnect |
| Values | "Who is asking?" | User ID, role from JWT |

### Why Context for User Data

- **Global variable:** Multiple concurrent requests overwrite each other
- **Function parameter:** Every handler signature must change
- **Context:** Each request carries its own isolated data

### Storing and Retrieving

```go
// Middleware stores:
ctx := context.WithValue(r.Context(), "user", claims)
next.ServeHTTP(w, r.WithContext(ctx))

// Handler retrieves:
claims := r.Context().Value("user").(*auth.Claims)
```

### Key Collision Risk

String keys like `"user"` can collide with other packages. Go convention:

```go
type contextKey string
const userClaimsKey contextKey = "user"
```

Private type = only your package can use this key.

---

## Chi Router: Use vs With

| Method | Behavior | Use Case |
|--------|----------|----------|
| `r.Use(mw)` | Mutates router, applies to all subsequent routes | Inside a Group |
| `r.With(mw)` | Returns new router, original unchanged | Single route protection |

### Common Mistake

```go
r.Group(func(r chi.Router) {
    r.With(middleware.Auth)  // Returns new router, discarded!
    r.Post("/movies", ...)   // Uses original r, unprotected
})
```

### Correct Usage

```go
r.Group(func(r chi.Router) {
    r.Use(middleware.Auth)   // Mutates this router
    r.Post("/movies", ...)   // Now protected
})

// Or for single route:
r.With(middleware.Auth).Post("/movies", handler)
```

---

## Import Aliasing

When two packages have the same name:

```go
import (
    chimw "github.com/go-chi/chi/v5/middleware"  // alias
    "github.com/yourproject/internal/middleware"
)

// Use:
r.Use(chimw.Logger)
r.Use(middleware.Authenticate)
```

---

## Token Transmission

### Where Tokens Live in HTTP

| Location | Problem |
|----------|---------|
| Query param `?token=...` | Logged in server logs, browser history, referrer headers |
| Request body | GET requests don't have bodies; must parse body before auth |
| **Header** | Metadata about request, not part of resource data |

### The Standard

```
Authorization: Bearer <token>
```

"Bearer" = whoever bears (carries) this token is authorized (RFC 6750).

---

## Refresh Tokens — Full Flow

### The Complete Lifecycle

```
USER LOGIN
==========
Client sends: { username, password }

Server does:
  1. Validate credentials
  2. Generate access token (JWT, 15min)
  3. Generate refresh token (UUID, stored in Redis with 7d TTL)
  4. Return BOTH tokens to client

Client stores both tokens locally.


NORMAL API USAGE
================
Client sends: GET /movies
              Authorization: Bearer <access_token>

Server does:
  1. Middleware extracts token
  2. Validates signature + expiration
  3. Passes request to handler

No Redis hit. No DB hit. Stateless.


ACCESS TOKEN EXPIRES (15 min later)
====================================
Client sends: GET /movies
              Authorization: Bearer <expired_access_token>

Server responds: 401 Unauthorized

Client thinks: "My access token is dead. I still have my refresh token."


TOKEN REFRESH
=============
Client sends: POST /auth/refresh
              Body: { "refresh_token": "<the-uuid>" }

Server does:
  1. Look up refresh_token in Redis
  2. Not found? → 401 (token revoked or expired)
  3. Found? → Get user_id from Redis value
  4. Delete OLD refresh token from Redis
  5. Generate NEW access token (JWT, 15min)
  6. Generate NEW refresh token (UUID, store in Redis, 7d TTL)
  7. Return both new tokens

Why delete + reissue? → "Token rotation" (see below)


LOGOUT
======
Client sends: POST /auth/logout
              Body: { "refresh_token": "<the-uuid>" }

Server does:
  1. Delete refresh token from Redis
  2. Done.

Access token still works until it expires (15min max).
That's the accepted tradeoff of stateless tokens.
```

### The Two Tokens Side by Side

```
                  Access Token          Refresh Token
Format:           JWT (signed)          UUID (random string)
Lives:            Client only           Client + Redis
TTL:              15 minutes            7 days
Used for:         API requests          Getting new access token
Validated by:     Signature check       Redis lookup
Revocable?:       No                    Yes (delete from Redis)
```

### Why Refresh Token is UUID, Not JWT

JWT's advantage is stateless validation. But refresh tokens are looked up in Redis anyway — you're already hitting a store. The self-contained nature of JWT adds no value here. A random UUID is simpler: look it up in Redis, get user_id, done. No signing, no claims parsing, no token-side expiration logic. Redis TTL handles expiration.

### Token Rotation (Delete Old, Issue New)

If attacker steals a refresh token and you DON'T rotate:
- Attacker uses it forever (7 days)
- You never know

If you DO rotate:
- Attacker uses stolen token → gets new pair, old token deleted
- Real user tries to refresh with old token → fails
- That failure is your **signal** that something is wrong

### Why Not Postgres for Refresh Tokens?

Refresh tokens are short-lived, high-churn data (created, looked up, deleted, expire). Redis handles TTL natively — expired keys disappear automatically. Postgres would need a background job to clean up expired rows. Postgres is for data you want to keep. Redis is for data you want to expire.

### The Stateless Tradeoff

After logout, the access token is still valid until it expires (15min). In that window, anyone holding it can use it.

**Why?** The password only matters at login. After that, the token IS the proof of identity. Every subsequent request sends ONLY the token — no username, no password:

```
POST /movies
Authorization: Bearer eyJhbGci...

{"title": "Inception"}
```

If an attacker gets the token (XSS, network sniffing, leaked logs), they send requests as you. The server validates the signature, sees it's not expired, and executes. It has no way to know it's not the real user.

**The password is the door. The token is the keycard inside the building.** Clone the keycard → walk around freely until it expires.

### The Revocation Problem

After logout, we delete the refresh token from Redis. But the access token? Still valid. JWT is stateless — there's no server-side record to delete. It's a monster that dies by itself.

**Can we fix this?** Yes — with a token blacklist in Redis:

```
Request → validate JWT signature → check Redis blacklist → proceed or reject
```

But at that point, you're hitting Redis on every request. The stateless advantage of JWT is gone. **You've reinvented sessions.**

### The Spectrum

```
Pure JWT (stateless)     → Fast, scalable, can't revoke
JWT + blacklist          → Revocable, but hitting Redis every request
Sessions (stateful)      → Full control, shared state from the start
```

There is no solution that gives you everything. Every system picks a point on this spectrum based on what damage looks like in their context.

Movie API? Pure JWT is fine. 15-minute blast radius, low-value targets.
Banking API? Never accept an unrevocable 15-minute window.

**Senior insight:** JWT is not "better" than sessions. It's a different set of tradeoffs. The right choice depends on what you're protecting.

---

## Implementation Gotchas (Learned the Hard Way)

### Error Handling: Log vs Return

```go
// WRONG — logs error, continues, panics later
user, err := cache.GetUserByRefreshToken(ctx, rdb, id)
if err != nil {
    log.Printf("failed: %v", err)  // logged but not returned
}
accessToken := auth.GenerateAccessToken(user.UserId, ...)  // user is nil → panic
```

**Rule:** If an error means you can't continue, **return immediately**. Logging without returning means "this error is not fatal, I have a fallback." If you don't have a fallback, return.

### Nil Pointer After Cache Miss

`GetUserByRefreshToken` returns `(nil, nil)` on cache miss — not an error, just "not found." If you don't check for nil before using the result, you get a nil pointer dereference.

```go
// Always check both
if err != nil { return err }
if user == nil { return errors.New("not found") }
```

### Type Assertion: Comma-OK Pattern

```go
// WRONG — panics if value is nil or wrong type
claims := r.Context().Value("user").(*auth.Claims)

// CORRECT — safe, no panic
claims, ok := r.Context().Value("user").(*auth.Claims)
if !ok { ... }
```

### Middleware Must Not Write Success Status

```go
// WRONG — writes 200 before handler runs, handler's status is silently dropped
w.WriteHeader(http.StatusOK)
next.ServeHTTP(w, r)

// CORRECT — middleware either rejects (write error + return) or passes through silently
next.ServeHTTP(w, r)
```

HTTP only allows one status code per response. If middleware writes 200, the handler's 201/204 is lost.

### r.Use() vs r.With() in Chi

```go
// WRONG — r.With() returns new router, discarded
r.Group(func(r chi.Router) {
    r.With(middleware.Auth)   // return value ignored!
    r.Post("/movies", ...)    // unprotected
})

// CORRECT — r.Use() mutates the router
r.Group(func(r chi.Router) {
    r.Use(middleware.Auth)
    r.Post("/movies", ...)    // protected
})
```

### Orphan Refresh Tokens

Each login creates a new refresh token without deleting the previous one. Multiple logins = multiple valid refresh tokens for the same user sitting in Redis.

**Production consideration:** Store a key like `user:<id>:refresh_token` that always holds the current token. New login overwrites it. Only one valid refresh token per user.

### Sessions: Why "Shared State"

Sessions are stored server-side. With multiple servers behind a load balancer:

```
                    ┌── Server A (has session)
User → Load Balancer├── Server B (no idea who this user is)
                    └── Server C (no idea who this user is)
```

Solutions (all add complexity): sticky sessions, shared Redis store, session replication. JWT avoids this — every server validates independently with just the secret key. But you lose instant revocation.

---

> this lab is very interesting.
