### 🔹 Problem 1: In-Memory User Store

Create a package:

internal/store

Inside it, create:

```go
type UserStore struct {
    users map[string]User
}
```
Tasks

- Initialize the store
- Add method: Create(user User)
- Add method: GetByID(id string)
- Add method: List() []User

---

### 🔹 Problem 2: Concurrency Safety (Important ⚠️)

Now imagine:

100 requests hitting /users at same time

Task

- Make your UserStore concurrency-safe (use `sync.Mutex`)

---

### 🔹 Problem 3: REST Endpoints

Implement these endpoints:

|Method	|Route	|Description|
|-------|-------|-----------|
|POST	|/users	|Create user|
|GET	|/users	|List users |
|GET	|/users/{id}	|Get user by ID |

Rules

- Use only net/http
- Parse {id} manually from URL
- Return proper HTTP status codes

---

### 🔹 Problem 4: Error Responses (Consistency)

Create a standard error response:

```go
{
  "error": "message"
}
```

Task

- Use it everywhere
- Never return raw text errors

### 🔹 Problem 5: Dependency Injection (Very Important Concept)

Refactor your handlers so that:

- They don’t create the store themselves
- Store is injected into handlers

Example idea:

```go
type UserHandler struct {
    store *UserStore
}
```

🧠 Learn why global variables are bad

---

### 🔹 Problem 6: Server Composition

Your main.go should:

- Create store
- Create handlers
- Attach routes
- Start server

No business logic inside main.