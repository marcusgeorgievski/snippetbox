# Processing Forms

- Parse and access form data sent in a `POST` request
- Validation checks on form data
- validation failures and repopulating form fields with previously submitted data
- Keep handlers clean by using helpers

## Parsing Form Data

1. Use `r.ParseForm()` to parse request body. Checks that body is well-formed, then stores form data in the `r.PostForm` map.
2. Use `r.PostForm.Get()` to retrieve value of a field. If no matching field, an empty string will be returned

- `r.ParseForm()` - Parse request body
- `r.PostForm` - Map of field names and values
- `r.PostForm.Get()` - Get a field name's value

Others:

- `r.PostFormValue()` - Shortcut, calls `r.ParseForm()` and fetches appropriate field value (ignores errors returned by `r.Parseorm(), ignore`)
- `r.Form()` - Contains form data from any `POST` request **and** any query string params

**Multipe-value fields**

- `r.PostForm.Get()` will only get the first value for a form field. If a field has multiple values (checkboxes, etc), we need to iterate over the underlying type of `r.PostForm` - `url.Values`
- `url.Values` is a `map[string][]string`

```go
for i, item := range r.PostForm["items"] {
    fmt.Fprintf(w, "%d: Item %s\n", i, item)
}
```

**Limiting form size**

- `POST` forms have a size limit of 10MB by default
  - Exception if form has the `enctype="multipart/form-data"` attribute and is setting multpart data, in which case there is no default limit
- `http.MaxBytesReader()` can set a specified limit

```go
r.Body = http.MaxBytesReader(w, r.Body, 4096)
```

**Query string parameters**

If a form submits data using a `GET` method, rather than `POST`, the form data will be included the URL query string parameter

```html
<form action="/foo/bar" method="GET">
  <input type="text" name="title" />
  <input type="text" name="content" />

  <input type="submit" value="Submit" />
</form>
```

Will result in

```
/foo/bar?title=value&content=value
```

These values can be retrievd with `r.URL.Query().Get("...")`

- Could also use `r.Form()` object

## Validating Form Data

https://www.alexedwards.net/blog/validation-snippets-for-go

## Displaying errors and repopulating fields

1. Add form data to `templateData` struct for form re-population
2. Create `snippetCreateForm` struct to hold form values and field errors
3. Validate, add to `FieldErrors` if validation check fails
4. If field errors, render create form again, and pass in `templateData` with `snippetCreateForm` on Form
5. Modify markup to accept input values if not nil or errors
6. Also consider default values in the create form handler

Why we don't use more 'RESTful' route structure like:

- `POST /snippets`
- `GET /snippets/{id}`

1. `GET /snippet/{id}` can conflict with `GET /snippet/create`. Although `id` will always be an in _in this case_, overlapping routes can be a source of bugs
2. If we submit the form at `GET /snippet/create` to `POST /snippets` and a re-render occurs to show a validation error, the url will change to `/snippets` which is confusing in terms of UX (especially if a `GET /snippets` does nothing/something else)

## Creating validation helpers

- Validator struct
- Embedded structs

- [Generics](https://go.dev/doc/tutorial/generics)

## Automatic form parsing

1. Form decoder
2. Struct tags on form struct
3. Decode in handler, check error (pass in form struct pointer)

- Type conversions are handled automatically too, such as `expires` being mapped to an int
- Check for invalid decoder error (not client's fault)
