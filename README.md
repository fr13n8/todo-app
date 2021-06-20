# REST API TODO in Go

## Development tools and principles
- Go Web Applications following the REST API Design
- Working with the framework <a href="https://github.com/gin-gonic/gin">gin-gonic/gin</a>
- The Clean Architecture approach to building an application structure. Dependency injection technique
- Working with Postgres DB. Running from Docker. Generation of migration files
- Application configuration using the library <a href="https://github.com/spf13/viper">spf13/viper</a>. Working with environment variables
- Working with BD using the library: <a href="https://github.com/jmoiron/sqlx">sqlx</a>
- Auth with JWT(get and refresh) and middleware
- Write SQL queries
- Graceful Shutdown

### DEMO
<a href="https://go-todo-app-backend.herokuapp.com/swagger/index.html">todo-app</a>

### How to run

```properties
make build && make run
```

If app startet first time need to run migrations

```properties
make migrate
```

### All commands

- Build
```properties
make build
```
- Run
```properties
make run
```
- Run Tests
```properties
make test
```
- Migrate up
```properties
make migrate_up
```
- Migrate down
```properties
make migrate_down
```
- Generate swagger documentation
```properties
make swag
```
- Generate mocks
```properties
cd pkg/service && go generate
```