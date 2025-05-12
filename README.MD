# BlitzHTTP

`blitzhttp` is a lightweight, high-performance HTTP router for Go, designed for simplicity and speed. It provides a clean API for routing, middleware, and route groups, with optimizations like pre-composed middleware, static route compilation, and zero-allocation wildcards. It supports method-specific routes, parameterized routes, wildcards, catch-all routes, and nested route groups.

## Features

- **Method-specific routes**: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, and `ANY`.
- **Parameterized routes**: Match paths like `/users/:id`.
- **Wildcard routes**: Handle paths like `/files/*`.
- **Catch-all routes**: Match any path with `ANY("*")`.
- **Route groups**: Organize routes with prefixes and shared middlewares.
- **Nested route groups**: Support hierarchical routing.
- **Middleware**: Apply multiple middlewares globally, per route, or per group, with short-circuiting for rejected requests.
- **Auto-OPTIONS handling**: Responds to `OPTIONS` requests with allowed methods.
- **High performance**: O(1) static lookups, switch-based method dispatch, and minimal allocations.

## Installation

Add `blitzhttp` to your Go project:

```bash
go get github.com/yourusername/blitzhttp@v0.1.0
```
