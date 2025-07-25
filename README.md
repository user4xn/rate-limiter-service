# Rate Limiter Service

A simple Rate Limiter designed for internal service-to-service communication, matched on (payment gateway architecture). This limiter ensures controlled traffic across internal APIs using Redis atomic operations for concurrency safety.
This rate limiter is built not to be accessed by end-users directly but by internal services that act as middleware.

## 🗺️ Architecture Flow

Summary:
1. Client → Sends request to Web API service (e.g., initiate payment).
2. Middleware intercepts and queries Rate Limiter Service.
3.Service checks Redis:
- If config for client+route exists → use it.
- If not → set default (e.g., limit=100, window=60s).
4. Redis key format:
- `fixed-window:{client_id}:{route}:{window_timestamp}`
6. Uses atomic INCR + EXPIRE to manage concurrent hits.
7. If limit exceeded → response returns DENY but still wrapped in 200 OK (internal system compatibility).

## 🚀 Features

- 🧠 Fixed Window rate limiting per client and route
- 🛡️ Internal-only usage (client ID + API key headers)
- 🛠️ Configurable limits via API
- ⚡ Redis-backed for distributed safety
- 🔐 Safe for concurrent requests using Redis atomic operations (`INCR`, `EXPIRE`)
- 📦 Dockerized for local or production deployment
- ✅ Unit tested: includes burst simulation & config update tests

---

## 🧩 Default Configuration
```
_____________________________
| Property      | Value     |
|---------------|-----------|
| Max Requests  | `100`     |
| Time Window   | `60s`     |
_____________________________

```

Custom per-client limits can be set through the config API.

---

## 📌 API Endpoints

### 🔍 1. Check Rate Limit

**POST** `/api/v1/rate/fixed-window`

#### Example `curl`
```bash
curl --location --request POST 'http://{your_address}/api/v1/rate/fixed-window' \
--header 'Api-Key: {your_api_key}' \
--header 'X-Client-Id: {client_id}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "route": "/api/v1/transactions"
}'
```
###  ⚙️ 2. Set Client Configuration
**PUT** `/api/v1/rate/fixed-window/set`

#### Example `curl`
```bash
curl --location --request PUT 'http://{your_address}/api/v1/rate/fixed-window/set' \
--header 'Api-Key: {your_api_key}' \
--header 'X-Client-Id: {client_id}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "route": "/api/v1/transactions",
    "limit": 5,
    "window": 20
}'

```
---
### 🧪 How To Run
Prerequisites:
- Docker & Docker Compose 

Clone the project
```bash
git clone https://github.com/user4xn/rate-limiter-service.git
```
Then access the project `dir` in `rate-limiter-service` folder

Setup .env file
```bash
cp .env.example .env
```
configure it match to your machine

Start the service
```bash
docker-compose up -d --build
```

The service will run at:
`http://localhost:8080`

Or if you want start service without docker

```bash
go run main.go
```

Require `go` and `redis` installed on your machine


---
### 🧪 Running Tests
Make sure have `go` installed on your machine
```bash
go run main.go test
```
Tests include:

- Basic rate limit validation
- Configuration updates
- Concurrent burst simulation
