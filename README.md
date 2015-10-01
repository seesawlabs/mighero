# mighero
SQL Migration Tool that uses our config files.

- [ ] Add ```migration_dir``` to the default.yml (Golang-api)
- [ ] Add ```driver``` to the default.yml (Golang-api)
```yaml
default:
    db: &db_default
        ip: localhost:3306
        user: 
        password: 
        name: kountable
        migration_dir: ./migrations
        driver: mysql
```
