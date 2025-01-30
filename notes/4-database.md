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
