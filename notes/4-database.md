# Database-driven responses

```sh
brew services start mysql
mysql -u root -p # empty password
```

```sql
CREATE USER 'web'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';
ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';
USE snippetbox;
```

- `go get` command will recursively download dependencies that the packages has

## Modules and Reproducible Builds

- `go.sum` - cryptographic checksums representing the content of the required packages
- `go mod verify` - verifies checksums of downloaded packages on my machine match entries in `go.sum`
- `go mod download` - download all dependencies at exact versions
- `go get -u <package>` - upgrade to latest available minor or patch release of an existing package. They are fixed otherwise. Use @x.x.x for specific version
- `go get <package>@none` - remove unused package
- `go mod tidy` - automatically remove unused packages from `go.mod` and `go.sum`

## Database Connection

- **Data Source Name (DSN)** - describes how to connect to a database
- `parseTime=true` - drive-specific param which instructs drivere to convert SQL `TIME` and `DATE` fields to Go `time.Time` objects
- `sql.Open()` - returns `sql.DB` object. It is a pool of many connection. Go manages the connections in this pool as needed, automatically opening and closing connections to the database via the driver.
- imports prefixed iwth an udnerscore don't use anything from the package, but they need the package's `init()` function to run. In this case, the `mysql` package needs to register itself with the `database/sql` package

## DB Functions

- `DB.Query()` is used for SELECT queries which return multiple rows.
- `DB.QueryRow()` is used for SELECT queries which return a single row.
- `DB.Exec()` is used for statements which don’t return rows (like INSERT and DELETE).
- `sql.Result.LastInsertId()` is used to get the id of a newly inserted record, not supported by postgres
- `sql.Result.RowsAffected()` number of rows affected

**Placehodler Parameters and `DB.Exec()`**

`?` inidicates a _placeholder parameter_ for MySql. This is used to contruct our query to help avoid SQL injection attacks from untrusted user-provided input. (postgres uses `$N`).

`DB.Exec()` behind the scenes:

1. Creates a _prepared statement_ on db using provided statement. DB parses and compiles the statement, then stores it ready for execution.
2. Passes the param values to the db. Db executes the prepared statement with params. Since the statement has already been compiles, the db treats the params as _pure data_, they cannot change the intent of the statement.
3. Closes (deallocates the prepared statement on the db)

## Single-record SQL Queries

- `rows.Scan()` - driver will automatically convert raw output from db to required native Go types. As long as the maped types are sensible, the converstions should work.

**Some Mappings**

- `CHAR`, `VARCHAR` and `TEXT` map to `string`
- `BOOLEAN` maps to `bool`
- `INT` maps to `int`; `BIGINT` maps to `int64`
- `DECIMAL` and `NUMERIC` map to `float`
- `TIME`, `DATE` and `TIMESTAMP` map to `time.Time` (remember `parseTime=true` in DSN, otherwise it returns a []byte)

**Specific Errors**

- `errors.Is(err1, err2)`- check if `err1` matches `err2`
- `err1 == err2` was the idomatic ways prior to 1.13. Go introduces wrapping errors which creates an entirely new error, so it is not possible to check the value of the original underlying error. `errors.Is()` unwraps errors before checking for a match.
- `errors.As()` Check if a potentially wrapped error has a specific type.

## Multiple-record SQL Queries

1. `rows, err := m.DB.Query(stmt)`
   - Query, then check for error
2. `defer rows.Close()`
   - Critical to ensure `sql.Rows` resultset is always properly closed and free up connections
3. `var snippets []Snippet`
   - Create slice of model
4. `for rows.Next() { ... }`
   - Iterate over resultset, prepares first (then each subsequent) row. If the iteration over all rows completes, the resultset closes itself and frees-up underlying db connection.
   - Create single model ovject, scan, check for error, append single row to slice
5. `if err = rows.Err(); err != nil { return nil, err }`
   - Retrieve any error that occured during iteration - don't assume successful iteration
6. `return snippets, nil`

## Transactions and other details

- Go does not handle `NULL` values well and will return an error if tried to convert to a string
- Transactions are useful to execute multiple SQL statements as a single atomic action
- Since `Exec()`, `Query()`, `QueryRow()` prepare statements at run time, a better idea could be to use `DB.Prepare()`. Thils will create the prepared statement once, and resuse it. This ccan be done in a function that returns a new model with a db connection. These prepared statements bind to specific connections in a pool, requiring re-preparation under high load, which can increase complexity and hit server limits. Using `Query()`, `QueryRow()`, and `Exec()` without explicit preparation is often simpler and sufficient.
