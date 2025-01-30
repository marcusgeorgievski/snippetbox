# Configuration and Error Handling

## Command-line Flags

- Get command-ling flags for different types with `flag.String()`, `Bool()`, `Int()`, etc.
- Can also use `flag.StringVar()`, `BoolVar()`, etc. if you want tostore in a struct

```go
type config struct {
    staticDir string
}

...

flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
```

Usage and information:

```sh
# Run with a flag
go run ./cmd/web -addr=:4000

# Get info on flags and usage
go run ./cmd/web -help
```

Can also store config in environment variables an access with `os.Getenv('foo')`. The drawback is that you cannot specify default values, get usage information, and it is always a string. Get the best of both worlds by passing env variables as flags.

```sh
export SNIPPETBOX_ADDR=":9999"
go run ./cmd/web -addr=$SNIPPETBOX_ADDR
```

## Structured Logging

- `log/slog`
- `slog.HandlerOptions` to define behaviour
- No equivalent to `log.Fatal()`, must call `os.Exit(1)` to terminate
- `slog.String(key, value)`, `slog.Any(key, value)`, etc.

## Dependency Injection

- [Some ideas](https://www.alexedwards.net/blog/organising-database-access)
