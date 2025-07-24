# ⚙️ How It Works

This service implements a **Fixed Window Rate Limiting** algorithm to control how many requests a client can make to a specific route within a defined time window (e.g., 100 requests per minute).

---

## 📊 Fixed Window Algorithm

### ✅ As Required Task

The **Fixed Window** algorithm divides time into equal-length intervals (windows), such as every 60 seconds. Each client+route combination gets its own counter per window.

If the number of requests exceeds the limit within the current window, further requests are **denied (marked as not allowed)** until the next window begins.

---

### 🧠 How It Works Internally

1. **Incoming Request**  
   The internal service makes a request to the rate limiter API with:
   - `client_id`
   - `route`
   - `API-Key`

2. **Key Generation**  
   A unique key is formed using the combination of:
`{client_id}:{route}:{current_window_timestamp}`
For example: `user123:/api/v1/order:2025-07-24T12:00:00`

3. **Counter Lookup in Redis**  
- If the key exists, the counter is incremented.
- If the key does **not** exist, it is created with a TTL equal to the window duration (e.g., 60s).

4. **Decision Logic**  
- If the counter **≤ max allowed requests** → allow.
- If the counter **> max allowed requests** → deny.

---

## 🛠 Example Flow
```
| Time       | Request # | Counter | Window | Result   |
|------------|-----------|---------|--------|----------|
| 12:00:01   | 1         | 1       | 12:00  | ✅ Allow |
| 12:00:20   | 2         | 2       | 12:00  | ✅ Allow |
| 12:00:59   | 101       | 101     | 12:00  | ❌ Deny  |
| 12:01:00   | 1 (new)   | 1       | 12:01  | ✅ Allow |
```
> A client can make up to 100 requests in each 60-second window.

---

## 🧠 Default vs Custom Configuration

- If a custom config for a `client_id+route` exists in Redis, it is used.
- If not, a **default config** is applied (e.g., 100 requests per 60 seconds).

---

## 🌐 Distributed Support (Optional)

- When using **Redis**, multiple instances of the rate limiter can share request counters.
- Redis operations (`INCR`, `EXPIRE`) ensure atomic and consistent behavior across distributed nodes.

---

## 🔐 Response Format

The rate limiter always responds with `200 OK`, and the actual decision is in the body:
```json
{
    "meta": {
        "message": "success",
        "code": 200,
        "status": "ok"
    },
    "data": {
        "status": "Allow", //OR Deny
        "limit": 100,
        "remain": 99,
        "reset_in_second": 41
    }
}
```
