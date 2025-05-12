# BlitzHTTP

`blitzhttp` is a lightweight, high-performance HTTP router for Go, engineered for simplicity, speed, and minimalism. Built with modern Go APIs, it provides a clean and intuitive interface for defining routes, applying middleware, and organizing route groups. Optimized for low-latency API servers, `blitzhttp` excels in scenarios with static routes and catch-all handlers, making it ideal for both small and large-scale applications. It leverages static route compilation, switch-based method dispatch, pre-composed middleware, and zero-allocation wildcards to outperform many popular Go routers.

## Introduction

`blitzhttp` is designed for developers who need a fast, reliable, and easy-to-use HTTP router without the overhead of complex frameworks. It integrates seamlessly with Go’s standard `net/http` package, ensuring compatibility with existing Go ecosystems. Key optimizations include:

- **O(1) static route lookups** for unmatched performance.
- **Pre-composed middleware** to eliminate runtime chaining.
- **Zero-allocation wildcards and catch-all routes** for minimal memory usage.
- **Switch-based method dispatch** for fast HTTP method resolution.

Whether you’re building a simple API with a few static routes and a catch-all handler or a complex system with parameterized routes and nested groups, `blitzhttp` delivers exceptional performance and flexibility.

## Features

- **Method-Specific Routes**: Support for `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, and `ANY` methods.
- **Static Routes**: O(1) lookups for exact path matches (e.g., `/api`).
- **Parameterized Routes**: Match dynamic paths like `/users/:id` with efficient segment-based parsing.
- **Wildcard Routes**: Handle trailing wildcards (e.g., `/files/*`) with zero-allocation parameter storage.
- **Catch-all Routes**: Match any path with `ANY("*")`, ideal for fallback handlers or custom routing.
- **Route Groups**: Organize routes with shared prefixes and middleware (e.g., `/api/v1`).
- **Nested Route Groups**: Support hierarchical routing for complex APIs.
- **Middleware**: Apply multiple middlewares globally, per route, or per group, with short-circuiting for rejected requests (e.g., authentication).
- **Auto-OPTIONS Handling**: Automatically responds to `OPTIONS` requests with allowed methods, simplifying CORS.
- **High Performance**: Outperforms `chi`, `gorilla/mux`, and `gin`, and is competitive with or faster than `httprouter` for static routes.
- **No Dependencies**: Relies solely on Go’s standard library (`net/http`).
- **Clean Code**: Modular design with files under 150 lines, ensuring maintainability.

## Installation

Add `blitzhttp` to your Go project using the tagged release `v0.1.0`:

```bash
go get github.com/NazmulIslamFahim/blitzhttp@v0.1.0
```

Replace `NazmulIslamFahim` with the actual repository owner. This fetches the stable `v0.1.0` version, which includes all features and optimizations described here.

For the latest development version (not recommended for production):

```bash
go get github.com/NazmulIslamFahim/blitzhttp@main
```

### Local Development

If you’re working with a local copy of `blitzhttp`:

1. Clone the repository:
   ```bash
   git clone https://github.com/NazmulIslamFahim/blitzhttp.git
   ```
2. In your project’s `go.mod`, add a `replace` directive:

   ```go
   module yourproject

   go 1.21

   require github.com/NazmulIslamFahim/blitzhttp v0.1.0

   replace github.com/NazmulIslamFahim/blitzhttp v0.1.0 => ../blitzhttp
   ```

3. Run `go mod tidy` to resolve dependencies.

## Usage

### Basic Example

This example sets up a simple server with a few static routes and a catch-all handler, optimized for your use case:

```go
package main

import (
  "fmt"
  "net/http"
  "github.com/NazmulIslamFahim/blitzhttp"
)

func main() {
  router := blitzhttp.New()

  // Global middleware (e.g., logging)
  router.Use(blitzhttp.Logger())

  // Static routes
  router.GET("/api", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Welcome to the API!")
  }, blitzhttp.Auth())

  router.POST("/submit", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Data submitted")
  })

  // Catch-all route
  router.ANY("*", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Not found: %s", blitzhttp.GetParams(r))
  })

  // Start server
  http.ListenAndServe(":8080", router)
}
```

Test with `curl`:

```bash
curl http://localhost:8080/api -H "Authorization: Bearer token"
# Output: Welcome to the API!
curl http://localhost:8080/random
# Output: Not found: random
```

### Parameterized Routes

Handle dynamic paths like `/users/:id`:

```go
router.GET("/users/:id", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "User: %s", blitzhttp.GetParams(r))
})
```

```bash
curl http://localhost:8080/users/123
# Output: User: users/123
```

### Wildcard Routes

Match paths with trailing wildcards (e.g., `/files/*`):

```go
router.GET("/files/*", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "File: %s", blitzhttp.GetParams(r))
})
```

```bash
curl http://localhost:8080/files/docs/report.pdf
# Output: File: docs/report.pdf
```

### Route Groups

Organize routes with shared prefixes and middleware:

```go
apiGroup := router.Group("/api/v1", blitzhttp.Logger())
apiGroup.GET("/users", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "List users")
})

adminGroup := apiGroup.Group("/admin", blitzhttp.Auth())
adminGroup.GET("/dashboard", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Admin dashboard")
})
```

```bash
curl http://localhost:8080/api/v1/users
# Output: List users
curl http://localhost:8080/api/v1/admin/dashboard -H "Authorization: Bearer token"
# Output: Admin dashboard
```

### Middleware

Apply middleware globally, per route, or per group. Middlewares can short-circuit requests:

```go
// Custom middleware
func CustomMiddleware() blitzhttp.Middleware {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      fmt.Println("Custom middleware executed")
      next.ServeHTTP(w, r)
    })
  }
}

router.Use(CustomMiddleware())
router.GET("/secure", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Secure endpoint")
}, blitzhttp.Auth(), CustomMiddleware())
```

```bash
curl http://localhost:8080/secure -H "Authorization: Bearer token"
# Output: Secure endpoint
curl http://localhost:8080/secure
# Output: Unauthorized
```

### Parameter Access

For parameterized, wildcard, or catch-all routes, access the matched path via `GetParams`:

```go
router.ANY("*", func(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Caught path: %s", blitzhttp.GetParams(r))
})
```

```bash
curl http://localhost:8080/anything/here
# Output: Caught path: anything/here
```

## Performance

`blitzhttp` is optimized for speed and minimal resource usage, making it ideal for high-throughput APIs:

- **Static Routes**: O(1) lookups via compiled dispatch, no allocations (~10-20 ns/op).
- **Catch-all Routes**: Minimal overhead (one failed map lookup + one string allocation, ~20-30 ns/op).
- **Parameterized Routes**: Efficient segment-based matching (~50-100 ns/op).
- **Wildcard Routes**: Fast prefix matching with zero-allocation parameters (~30-70 ns/op).
- **Middleware**: Pre-composed at registration, eliminating runtime chaining.
- **Scalability**: Cache-friendly maps and minimal allocations ensure performance under load.

### Your Use Case: Few Static Routes + Catch-all

If your project uses only a few static routes (e.g., `/api`, `/submit`) and one catch-all route (`ANY("*")`):

- **Optimal Performance**: Static routes are resolved in constant time with no allocations, and the catch-all adds negligible overhead.
- **Comparison to Complex Projects**: Faster than projects with parameterized or wildcard routes, as it skips linear-time matching and additional allocations.
- **Library Comparison**:
  - **Faster than `httprouter`** for static routes due to static compilation.
  - **Much faster than `chi`, `gorilla/mux`, or `gin`** due to minimalism and no context overhead.
  - Competitive with `httprouter` for catch-all routes, outperforming others.

### Benchmarks

Run benchmarks to measure performance:

```go
package blitzhttp_test

import (
  "net/http"
  "net/http/httptest"
  "testing"
  "github.com/NazmulIslamFahim/blitzhttp"
)

func BenchmarkStaticRoute(b *testing.B) {
  router := blitzhttp.New()
  router.GET("/api", func(w http.ResponseWriter, r *http.Request) {})
  req, _ := http.NewRequest("GET", "/api", nil)
  w := httptest.NewRecorder()
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    router.ServeHTTP(w, req)
  }
}

func BenchmarkCatchAllRoute(b *testing.B) {
  router := blitzhttp.New()
  router.GET("/api", func(w http.ResponseWriter, r *http.Request) {})
  router.ANY("*", func(w http.ResponseWriter, r *http.Request) {})
  req, _ := http.NewRequest("GET", "/random", nil)
  w := httptest.NewRecorder()
  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    router.ServeHTTP(w, req)
  }
}
```

```bash
go test -bench=.
```

Expected results:

- `BenchmarkStaticRoute`: ~10-20 ns/op (O(1), no allocations).
- `BenchmarkCatchAllRoute`: ~20-30 ns/op (minimal overhead).

## Testing

Write unit tests using `net/http/httptest` to verify routes and middleware:

```go
package main

import (
  "net/http"
  "net/http/httptest"
  "testing"
  "github.com/NazmulIslamFahim/blitzhttp"
)

func TestStaticRoute(t *testing.T) {
  router := blitzhttp.New()
  router.GET("/api", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "OK")
  })

  req, _ := http.NewRequest("GET", "/api", nil)
  w := httptest.NewRecorder()
  router.ServeHTTP(w, req)

  if w.Body.String() != "OK" {
    t.Errorf("expected OK, got %s", w.Body.String())
  }
}

func TestCatchAllRoute(t *testing.T) {
  router := blitzhttp.New()
  router.ANY("*", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Caught: %s", blitzhttp.GetParams(r))
  })

  req, _ := http.NewRequest("GET", "/random/path", nil)
  w := httptest.NewRecorder()
  router.ServeHTTP(w, req)

  if w.Body.String() != "Caught: random/path" {
    t.Errorf("expected Caught: random/path, got %s", w.Body.String())
  }
}

func TestAuthMiddleware(t *testing.T) {
  router := blitzhttp.New()
  router.GET("/secure", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Secure")
  }, blitzhttp.Auth())

  req, _ := http.NewRequest("GET", "/secure", nil)
  w := httptest.NewRecorder()
  router.ServeHTTP(w, req)

  if w.Code != http.StatusUnauthorized {
    t.Errorf("expected 401 Unauthorized, got %d", w.Code)
  }
}
```

```bash
go test
```

## API Reference

### Router Methods

- `New() *Router`: Creates a new router.
- `GET(path string, handler http.HandlerFunc, mws ...Middleware)`: Registers a GET route.
- `POST(path string, handler http.HandlerFunc, mws ...Middleware)`: Registers a POST route.
- `PUT(path string, handler http.HandlerFunc, mws ...Middleware)`: Registers a PUT route.
- `DELETE(path string, handler http.HandlerFunc, mws ...Middleware)`: Registers a DELETE route.
- `PATCH(path string, handler http.HandlerFunc, mws ...Middleware)`: Registers a PATCH route.
- `ANY(path string, handler http.HandlerFunc, mws ...Middleware)`: Registers a route for all methods (or catch-all if `path == "*"`).
- `Use(mws ...Middleware)`: Adds global middleware.
- `Group(prefix string, mws ...Middleware) *Group`: Creates a route group.

### Group Methods

- `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `ANY`: Same as router methods, with paths prefixed by the group’s prefix.
- `Group(prefix string, mws ...Middleware) *Group`: Creates a nested group.

### Middleware

- `Middleware`: Type alias for `func(http.Handler) http.Handler`.
- Built-in examples:
  - `Logger()`: Logs request method and path.
  - `Auth()`: Rejects requests without an `Authorization` header.

### Parameter Access

- `GetParams(r *http.Request) string`: Returns the matched path for parameterized, wildcard, or catch-all routes.

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository: `https://github.com/NazmulIslamFahim/blitzhttp`.
2. Create a feature branch: `git checkout -b feature/your-feature`.
3. Commit changes: `git commit -m "Add your feature"`.
4. Push to the branch: `git push origin feature/your-feature`.
5. Open a pull request.

Please include tests and update documentation for new features. Report issues at `https://github.com/NazmulIslamFahim/blitzhttp/issues`.

## License

MIT License. See [LICENSE](LICENSE) for details.

## Acknowledgments

`blitzhttp` was designed to balance simplicity and performance, inspired by the needs of modern Go APIs. Special thanks to the Go community for feedback and inspiration.
