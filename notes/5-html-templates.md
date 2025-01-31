# Dyamic HTML Templates

## Basics

- `{{ define "name" }} content {{ end }}`
  - Used for HTML, page titles
  - An html file may contain multiple; one for the html, one for the title to be used in a _base_ template
- `{{ template "name" . }}`
  - Invokes a template (often from another template)
- `{{ block "name" . }} default content {{ end }}`
  - Invoke a template with default content if it does not exist

## Dynamic Data

Pass data into the `data` parameter

```go
data := templateData{
    Snippet: snippet,
}

err = ts.ExecuteTemplate(w, "base", data)
```

Access it in the templates with

```
{{.Snippet.Title}}
```

**Notes**

- The `html/template` pkg escapes any data between `{{ }}` tags. This helps avoid XSS attacks
- If the type that you’re yielding between {{ }} tags has methods defined against it, you can call these methods (so long as they are exported and they return only a single value — or a single value and an error)
  - `<span>{{.Snippet.Created.Weekday}}</span>` would work since `Created` is a `time.Time` object
  - Parameters can be passed with spaces: `{{ ... .AddDate 0 6 0}}` (adds 6 months to a time)

## Template actions and functions

| Expression                              | Description                                                                                                                                                                                                                                              |
| --------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `{{if .Foo}} C1 {{else}} C2 {{end}}`    | If `.Foo` is not empty, render the content `C1`; otherwise, render the content `C2`.                                                                                                                                                                     |
| `{{with .Foo}} C1 {{else}} C2 {{end}}`  | If `.Foo` is not empty, set dot to the value of `.Foo` and render `C1`; otherwise, render `C2`.                                                                                                                                                          |
| `{{range .Foo}} C1 {{else}} C2 {{end}}` | If the length of `.Foo` is greater than zero, loop over each element, setting dot to the value of each element and rendering `C1`. If the length of `.Foo` is zero, render `C2`. The underlying type of `.Foo` must be an array, slice, map, or channel. |

- For the above, else clauses are optional

| Expression                     | Description                                                                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------- |
| `{{eq .Foo .Bar}}`             | Yields true if `.Foo` is equal to `.Bar`.                                                                                                |
| `{{ne .Foo .Bar}}`             | Yields true if `.Foo` is not equal to `.Bar`.                                                                                            |
| `{{not .Foo}}`                 | Yields the boolean negation of `.Foo`.                                                                                                   |
| `{{or .Foo .Bar}}`             | Yields `.Foo` if `.Foo` is not empty; otherwise yields `.Bar`.                                                                           |
| `{{index .Foo i}}`             | Yields the value of `.Foo` at index `i`. The underlying type of `.Foo` must be a map, slice, or array, and `i` must be an integer value. |
| `{{printf "%s-%s" .Foo .Bar}}` | Yields a formatted string containing the `.Foo` and `.Bar` values. Works in the same way as `fmt.Sprintf()`.                             |
| `{{len .Foo}}`                 | Yields the length of `.Foo` as an integer.                                                                                               |
| `{{$bar := len .Foo}}`         | Assigns the length of `.Foo` to the template variable `$bar`.                                                                            |

**Notes**

- Combine functions
  - `{{if (gt (len .Foo) 99)}} C1 {{end}}`
  - `{{if (and (eq .Foo 1) (le .Bar 20))}} C1 {{end}}`
- Control loop

```go
{{range .Foo}}
    // Skip this iteration if the .ID value equals 99.
    {{if eq .ID 99}}
        {{continue}} // Or break
    {{end}}
    // ...
{{end}}
```

## Caching Templates
