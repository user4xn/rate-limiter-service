## 📌 Assumptions and Limitations

This section outlines the assumptions made during development and the known limitations of the current version of the rate limiter.

---

### ✅ Assumptions

- The rate limiter is **only used internally** by trusted services within a secure network.
- Each internal service integrates with the rate limiter via HTTP middleware.
- Clients are uniquely identified by a `Client-ID` header or similar metadata.
- Each route can have its own rate limit configuration.
- Redis is available if distributed support is required; otherwise, the service can fall back to in-memory mode (less scalable).

---

### ⚠️ Limitations

1. **HTTP Latency Overhead**  
   - Because each request involves an **HTTP call to an external rate limiter**, there's a small but noticeable **latency penalty**, especially under high traffic.
   - 🔧 **Possible Improvement**: Switching to **gRPC** would reduce overhead and increase performance due to persistent connections and binary serialization.

2. **Single Point of Failure with Redis**  
   - If Redis is used and becomes unavailable, rate limiting will fall back to default or in-memory behavior, which may be inconsistent across distributed instances.

3. **No Dynamic TTL Reset**  
   - The TTL (expiration) for a rate window is set when the counter is first created and does not reset with each request. This is expected behavior in fixed window models but may not suit bursty or session-based traffic.

4. **API Key Rotation**  
   - API Key are static and stored in configuration or memory. There is no automated rotation or revocation mechanism implemented yet, but it can be improved.

---

### 🧭 Possible Future Improvements

- Implement **gRPC support** for more efficient inter-service communication.
- Add **support for additional algorithms** (e.g., sliding window, token bucket).
- Provide **metrics and dashboards** for monitoring request rates and limiter performance.
- Enable **rate limiting per IP address** or other attributes beyond `Client-ID`.

---