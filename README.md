# GOCMS

An aspiring CMS and plugin-able module for variety of tasks. 

**AI Disclosure: More than 90% of total LOC of `**/*_test.go` files are written by AI. Rest of the code is written with almost zero (google's search also shows AI powered answers so taking into account that) help from AI in an attempt to valdidate if I can still write code like 2020.**


## Development (As we aren't anywhere prod ready yet)
- Start the dev server:
```bash
make dev
```

- (Install goose before running this) Run the [migrations](./migrations) once the DB is up:
```bash
make migrate
```

### New Migrations?
- Create the migration file using:
```bash
make migrate-create NAME="<a-good-intent-for-the-migration>"
```

- Edit the sql file created in [migrations](./migrations) folder


### Added a new package? Wanna rebuild the dev-server?
- Run the build command to trigger a docker rebuild
```bash
make dev-build
```

### Run Tests (Please write tests... you know you need it)
- Run the test command
```bash
make test
```

## Design Decisions:

### Database and Driver
- `sqlx` is used along with `lib/pq` for the postgres db interactions.
- using `goose` for running/creating/tracking db migrations
- Database initialized once and passed down to handlers, services and repositories. `sqlx.DB` instance
- Transactions are and need to be supported in Repo level functions that modify DB state (like `DELETE`, `UPDATE` or `INSERT` queries)

### Code Structure
- Need to adhere to practices that enable testability of the platform
- Repo, Service and Handler layer abstraction helps ensure functionality is provided by interfaces (making callers testable)
- Individual modules should document any relevant domain decisions made while developing. (Example: content_type slugs in the CMS module are stored in the DB prefixed with their account's id and are stripped back when responding back in fetch APIs. Functionality enabled and tested by domain level methods)
