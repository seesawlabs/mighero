# mighero - Migration Hero
SQL Migration Tool that uses our config files.

### Installation

``` go get github.com/kountable/mighero/cmd/mighero ```

### Usage

You need to be in the project root directory (mighero search for the config files in env/ folder) to execute the command.

you can change this behavior by specified the path to the default and env config file, using -def and -env flags.


```shell
$ mighero create InitDatabase
Version 2 migration files created in db/migrations:
0001_InitDatabase.up.sql
0001_InitDatabase.down.sql

$ tree db/migrations
db/migrations
├── 0001_InitDatabase.down.sql
├── 0001_InitDatabase.up.sql

$ mighero up                
> 0001_InitDatabase.up.sql

0.0677 seconds

$ mighero down
< 0001_InitDatabase.down.sql

0.0751 seconds

```


This project use the [Migrate](github.com/mattes/migrate) Lib.
