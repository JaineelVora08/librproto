
# LIBR Prototype

A minimal version of the LIBR system to simulate decentralized message moderation using Go, PostgreSQL, goroutines, and channels.

---

## 🚀 Features

- **Submit messages** via REST API (`/submit`)
- **Moderation** by 3 simulated goroutines (moderators)
- **Consensus logic**: Accept if ≥2 moderators approve within timeout
- **Persistent storage** using PostgreSQL
- **Retrieve messages** via `/fetch/{timestamp}` and `/fetchall`

---

## 🛠️ Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/JaineelVora08/librproto.git
cd librproto
```

### 2. Environment Variables

Create a `.env` file in the root directory with your PostgreSQL credentials:

```
connection_string=postgres://username:password@localhost:5432/libr
```

### 3. PostgreSQL Setup

Make sure PostgreSQL is running. Then create the database:

```sql
CREATE DATABASE libr;
```

The table will be auto-created if it doesn't exist.

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run the Server

```bash
go run main.go
```

The server listens on: `http://localhost:8000`

---

## 📡 API Endpoints

### ➕ `POST /submit`

Submits a message for moderation.

**Request:**

```json
{
  "content": "This is a test message."
}
```

**Response (if accepted):**

```json
{
  "id": "uuid",
  "timestamp": 1744219507,
  "status": "accepted"
}
```

**Response (if rejected):**

```json
{
  "id": "uuid",
  "timestamp": 1744219507,
  "status": "rejected"
}
```

---

### 📥 `GET /fetch/{timestamp}`

Fetches messages with `status = accepted` for a specific UNIX timestamp.

**Response:**

```json
[
  {
    "id": "uuid",
    "content": "This is a test message.",
    "timestamp": 1744219507,
    "status": "accepted"
  }
]
```

---

### 📋 `GET /fetchall`

Returns all accepted messages.

---

## 🔁 Moderation Process

- 3 moderators simulated using goroutines
- Each randomly approves or rejects a message
- Delay of **1–3 seconds** simulated using `time.Sleep`
- `context.WithTimeout` limits total moderation to **3 seconds**
- Approval requires **2 of 3 moderators to accept**

---

## 📷 Screenshots

### ✅ API Request and Response (via Postman)

![Postman Submit API](./screenshot/Screenshot%202025-06-14%20115549.png)


---

### 🗂️ Logs - Moderation & Storage

![Console Logs](./screenshot/Screenshot%202025-06-14%20115617.png)

---

### 📥 Fetching Accepted Messages

![Fetch API](./screenshot/Screenshot%202025-06-14%20120023.png)

---

## 🧱 Technologies Used

- **Go**
- **PostgreSQL**
- **Gorilla Mux** – for routing
- **PGX** – for DB interaction
- **UUID** – for unique message IDs
- **.env** – for secure credentials

---
