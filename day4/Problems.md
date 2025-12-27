### 🔹 Problem 1: Introduce Chi Router

Install Chi:

```bash
go get github.com/go-chi/chi/v5
```

Task

- Replace http.HandleFunc
- Use chi.NewRouter()
- Move all routes to router

### 🔹 Problem 2: Route Grouping (Very Real-World)

Group your user routes:

r.Route("/users", func(r chi.Router) {
    r.Post("/", ...)
    r.Get("/", ...)
    r.Get("/{id}", ...)
})


Task

- Clean, readable routes
- No duplicated /users


### 🔹 Problem 3: Middleware – Global vs Route-Level

Create:

- Global middleware:
- Logging
- Request ID
- Route-level middleware:
- Allow only POST for /users


### 🔹 Problem 4: Panic Recovery Middleware

- Create a middleware that:
- Recovers from panics
- Returns 500 Internal Server Error
- Logs panic details


### 🔹 Problem 5: Request Context Usage (Revisit Context Lightly)

Use context to:

- Pass request ID from middleware to handler
- Access it inside handler

⚠️ Context is request-scoped, not global state.

### 🔹 Problem 6: Clean Folder Structure

Refactor to something like:

```text
cmd/server/main.go
internal/handlers
internal/store
internal/middleware
internal/models
```

Task

- No circular imports
- Clean package boundaries
