# Middleware

> - Building custom middleware
> - HTTP headers
> - Log requests
> - Recover panics
> - Middleware chains

**Middleware Pattern**

```go
func middlware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Middleware logic here...
        next.ServeHTTP(w, r)
    })
}
```

**Positioning**

Positioning the middleware before the servemux. Executes on every request the servemux handles:

- `middleware -> servemux -> handler`

Position middlware after the servemux, runs on specific routes:

- `servemux -> middleware -> handler`

In each case, the chain will reverse when the handler returns

## Setting Common Headers

- `Content-Security-Policy`
- `Referrer-Policy`
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: deny`
- `X-XSS-Protection: 0`

## Panic Recovery

- Go's HTTP server assumes the effect of a panic is isolate to the goroutie serving the request - remember, each request is handled in its own goroutine
- A panic will result in an empty response on the client's side - not a great user experience
- We can fix this with a panic recovery middleware which calls `app.serverError()`

- If a handler spins up another goroutine, this panicc recovery middleware will not handle it. Must make sure those panics are handled separately

## Middleware Chains
