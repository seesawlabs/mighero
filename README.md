# mighero - Migration Hero
SQL Migration Tool that uses our config files.

### TO DO

- [x] Add ```migration_dir``` to the default.yml (Golang-api)
- [x] Add ```driver``` to the default.yml (Golang-api)

Example 

```yaml
default:
    ...
    db: &db_default
        ip: localhost:3306
        user: 
        password: 
        name: kountable
        migration_dir: ./migrations
        driver: mysql
    ...
```

Use the [Migrate](github.com/mattes/migrate) Lib.
