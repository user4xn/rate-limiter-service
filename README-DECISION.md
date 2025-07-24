## 💡 Design Decisions
This rate limiter was developed as part of a **technical interview assignment**. I designed it based on how **Doitpay**, a payment gateway system, handles traffic internally.

Since Doitpay operates as an internal payment platform and does not expose rate limiting directly to external clients, this rate limiter was intentionally built as an **internal service**, not a public-facing one.

---


### 🔒 Internal-Only Design

**Why**:  
To reflect a realistic internal architecture like Doitpay’s, where services communicate with each other securely behind the scenes.

**How**:
- This service acts as a **standalone internal rate limiter**, not as public middleware.
- Each internal service integrates with it via a **custom middleware** that sends rate-limit check requests.
- Requests must include:
  - `API-Key` for internal authentication.
  - A `Client-ID` (used for identifying the request source).
  - The target route being accessed.

---

### ⚙️ Architecture Highlights

- **Standalone HTTP Service**: Decouples rate limiting logic for reusability.
- **Redis (Optional)**: Enables distributed and scalable counter storage.
- **Default Config Fallback**: Ensures limits even if a client config is missing.
- **Always Returns 200 OK**: Keeps internal systems from breaking, even when over the limit.
- **Middleware Handles Enforcement**: Internal services decide whether to block or proceed, based on limiter’s `allowed: true/false` response.

---

### 🧪 Built for Internal Control & Safety

This setup is optimized for **internal usage**, where the goal is not just to block traffic, but to **track and control it in a centralized, safe, and flexible way** — especially in high-traffic systems like payment gateways.

---
Here's an overview of the decisions behind the architecture and how each part works.

### 🧱 1. Standalone Service Architecture

**Why**:  
To decouple rate limiting from business logic, and allow **multiple internal services** to integrate with a single limiter endpoint.

**How**:  
Each internal service sends a request to the limiter before processing a client request. The limiter returns a decision (`allowed: true/false`) based on the client ID and route.

---

### ⏱️ 2. Fixed Window Rate Limiting

**Why**:  
Chosen for its **simplicity and performance**. It’s sufficient for most internal use cases where evenly distributed traffic is expected.

**How**:  
- Each client+route has a counter stored in memory or Redis.
- When a request comes in:
  - If a window key doesn’t exist, it’s created and expires in N seconds (e.g., 60).
  - The counter is incremented atomically.
  - If the counter exceeds the limit, the request is denied.

---

### 🧠 3. Redis

**Why**:  
To enable horizontal scaling and shared counters across multiple limiter instances.

**How**:  
- Uses Redis `INCR` and `EXPIRE` commands to manage request counts.
- Ensures atomicity and time-bound TTLs across nodes.
- Automatically falls back to default config if client config is not found.

---

### ⚙️ 4. Default Rate Limit Fallback

**Why**:  
To provide protection even for unregistered clients or routes.

**How**:  
- If no custom config is stored in Redis, the limiter uses a default setting (e.g., `100 RPM`).
- This prevents unconfigured clients from making unlimited requests.

---

### 🔐 5. Internal-Only Access with API Key

**Why**:  
To restrict access and ensure only **trusted internal services** can use the limiter.

**How**:  
- Each internal service must send a valid `API-Key` header.
- The limiter validates the key before processing the rate check.

---

### ✅ 6. Always Returns Success (200 OK)

**Why**:  
To maintain consistent communication flow between internal services, even when a request is over the limit.

**How**:  
- The limiter never breaks the chain with an HTTP error.
- Instead, it includes a flag in the response:  
  ```json
  { "allowed": false }
