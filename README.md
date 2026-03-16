# GOCMS

An aspiring CMS and plugin-able module for variety of tasks. 


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
