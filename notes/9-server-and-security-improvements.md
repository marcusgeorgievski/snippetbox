# Server and Security Improvements

## `http.Server` struct

It is more common in real world applications to manually create and use a `http.Server` struct instead. This open up the opportunity to customize the behaviour of the server.

```go
server := &http.Server{
    Addr: *addr,
    Handler: app.routes(),
}

logger.Info(fmt.Sprintf("starting server at http://localhost%s", *addr), slog.String("addr", *addr))
```

By default, it writes these entries using the standard logger — which means they will be written to the standard error stream (instead of standard out like our other log entries), and they won’t be in the same format as the rest of our nice structured log entries.

- Not possible to configure `http.Server` to use our structured logger. We instead need to convert our structured logger handler into a `*log.logger`
