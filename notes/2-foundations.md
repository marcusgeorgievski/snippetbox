# Foundations

>

**Go Modules**

Include a canonical name or identifier for your project which can be any string, but the important thing is uniqueness. Modules provide some advantages such as managing dependencies, avoiding [supply-chain attacks](https://go.dev/blog/supply-chain), and ensuring reproducible builds.

**Web Application Basics**

- **Handler** - Execute application logic and write HTTP response headers and bodies
- **Router (servemux)** - Stores mapping between URL routing patterns and corresponding handlers
- **Server** - Build in web server to listen for incoming requests

**Trailing Slash**

- Patterns without trailing slashes will only be matched when request URL path matches exactly in full
- Routes that end in a trailing slash like `/static/`are known as a _subtree path pattern_. They are matches whenever the start of a request URL path matches the subtee path, like `/static/**`
- To prevent _subtree path pattern_ from acting like they have a `**` wildcard at the end, append `{$}` to the end of the path, which effectively means match a single slash and nothing else. Useful for `/` route. It is only permitted at the end of a subtree path pattern.

**Additional Servemux Features**

- Request URL paths are automatically _sanitized_. If the request path contains any `.` or `..` elements or repeated slashes (`/foo/..//bar`), the user will automatically be redirected to an equivalent clean URL (`/foo/bar`).
- If a servemux is not registered, Go uses the default servemux and stores it in `http.DefaultServeMux`. This is not recommended:
  - Global variables mean any code any access it and register routes
  - Third-party codde can register routes to expose a malicious handler

**WildCard Route Patterns**

`foo/{user}/{post}`

- Each segment can only have one wildcard identifier, and it must fill the whole segment
  - `/date/{y}-{m}-{d}` or `/{slug}.html` are not valid
- Retrieve with `r.pathValue()`, returns a string
- If a request URL path can match two route patterns, _the most specific route pattern wins_
- **Remainder Wildcards:** `/foo/{bar...}` can match any number of segments after `/foo/` which can be accessed from `bar`

**Method-based routing**

- Methods such as `GET` are case sensitive and should always be uppercase, followed by at least one whitespace
- Including a method in the path pattern increases specificity

## Customizing Responses

**Status Codes**

- `w.WriteHeader()`
- Sets custom status code, can only be called once per rseponse
- If not called explicitly, first call to `w.Write()` will send a `200`
- `net/http` provides many constants for [HTTP status codes](https://pkg.go.dev/net/http#pkg-constants)

**Headers**

- `w.Header().Add()`
- Must be called **before** `w.WriteHeader()` or `w.Write()`. Any changes afterward do not effect the header the user sreceives
- `w.Header()` also has `Set() Del() Get() Values()` methods to manipulate and read from header map too
- Header name will always be canonicalized with `textproto.CanonicalMIMEHeaderKey()` which converts the first and any letter following a hyphen to uppercase, and the rest lowercase. Calling methods on the header name is case-insensitive (all lowercase for HTTP/2).
  - Use underlying map access to ignore this behaviour: `w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}`

**Bodies**

- Common to pass the `http.ResponseWrite` to another fn that writes the response for you
- Response writer satisfies the `io.Writer` interface

**Additional Info**

Go _sniffs_ the response body to automatically set the `Content-Type` header. If it can't guess, will fallback to `Content-Type: application/octet-stream`. Cannot distinguish JSON from plain text, so JSON responses will be sent with a `Content-Type: text/plain; charset=utf-8` header. Set properly with `w.Header().Set("Content-Type", "application/json")`

**Project Structure**

- [Server project](https://go.dev/doc/modules/layout#server-project)
- `cmd` - application-specific code
- `internal` - non-application-specific code. **Special Behaviour:** any packages which live under this directory can only be imported by code inside the parent of the `internal` directory (anything in `snippetbox` in our case). Cannot be imported by code outside our project
- `ui` - ui assets, html, css, images

## Templating

- `ts, err := template.ParseFiles()` - Parses html files to a template. Path must be relative to working dir, or absolute path of the project directory
- `ts.Execute(w, data)` - Writes template to a write (response body in our case)
- `ts.ExecuteTemplates(w, name, data)` - Write content of specific template, which will in turn invoke other template's data as needed

- `{{define "..."}}` - Define a distinct name template
- `{{template "..." .}}` - Denotes we want to invoke other named templates with the respective name at the given location. `.` represents dynamic data
- `{{block "..." .}}` - acts like `{{template}}`, but you can provide default content inside until `{{end}}`

```html
{{block "sidebar" .}}
<p>My default sidebar content</p>
{{end}}
```

`http.FileServer`

- Serves files over HTTP from a specifi directory
- Santizes request paths by running them through `path.Clean()`. Removes `.` or `..` to prevent directory traversal attacks
- Range requests supported
- `Content-Type` automatically set from `mime.TypeByExtension()`. Add your own type with `mime.AddExtenstionType()`
- Performance wise, Windows and Unix-based OSs will likely cache and serve from RAM rather than slow round-trip to disk
- Could also serve a single file with `http.ServeFile(w, r, "./ui/static/file.zip")`. This does not sanitize the path, so use `filepath.clean()` if constructing a file path from user input
- More about disabling directory listings

## http.Handler Interface

Handler is an object which satisfies the `http.Handler` interface, meaning it must have a `ServeHTTP()` method with the exact signature `ServeHTTP(http.ResponseWriter, *http.Request)`

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

A handler in its simplest form may look like this:

```go
type home struct {}

func (h *home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("This is my home page"))
}
```

Which can be registered with the `Handle` method

```go
mux.Handle("/", &home{})
```

When the servemux receives an HTTP request for `/`, it will then call the `ServeHTTP()` method of the `home` struct - which in turn writes the HTTP response

### Handler Functions

Since the above process is longer, its more common to write handlers as a normal function that takes in the same parameters. It does not have a `ServeHTTP()` method, so in itself it is not a handler. We can transform it to a handler using the `http.HandlerFunc()` adapter.

This automatically adds a `ServeHTTP()` method to the home function, which will then call the original function's code.

```go
mux.Handle("/", http.HandlerFunc(home))
```

`HandleFunc()` is syntactic sugar that transforms a function to a handler and registers it in one step.

```go
mux.HandlerFunc("/", home)
```

## Static Files

## Notes

- `/` treated as a catch-all if another path does not match
- `http.ListenAndServe()` expects network address in `host:port` format. Uses all of computer's available netwwork interfaces if omitted. Generally only need host.
- `go run` compiles code, created executable binary in `/tmp`, and executes it in one step. Expects a list of `.go` files or full module path. The following are equivalent in our case:
  - `go run .`
  - `go run main.go`
  - `go run github.com/marcusgeorgievski/snippetbox`
  - Requests are handled concurrently
