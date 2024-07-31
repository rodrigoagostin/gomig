
# GoMig

GoMig is a Golang library and CLI tool for generating and managing PostgreSQL database migrations. It allows you to create, alter, and manage database schemas with simple commands, inspired by Rails migrations.

## Installation

To use GoMig as a library in your Go project, add it to your module:

```sh
go get github.com/rodrigoagostin/gomig
```

To install the GoMig CLI tool, use the following command:

```sh
go install github.com/rodrigoagostin/gomig/cmd/gomig@latest
```

Ensure that your `GOPATH` is included in your `PATH` environment variable so that you can run the `gomig` command from any directory.

## Usage

### CLI Tool

You can generate different types of migrations using the `generate` command.

#### Create Table

To create a new table:

```sh
gomig generate create_<table_name> <column_name:type> <column_name:type> ...
```

Example:

```sh
gomig generate create_users "name:varchar(200)" "email:varchar(200)" "active:boolean"
```

This will generate a migration file in the `migrations` directory:

```sql
-- Migration: create_users
-- Created at: 20240101123456

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200),
    email VARCHAR(200),
    active BOOLEAN,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Alter Table

To alter an existing table by adding new columns:

```sh
gomig generate alter_<table_name> <column_name:type> <column_name:type> ...
```

Example:

```sh
gomig generate alter_users "age:int" "address:varchar(200)"
```

This will generate a migration file:

```sql
-- Migration: alter_users
-- Created at: 20240101123456

ALTER TABLE users
    ADD COLUMN age int,
    ADD COLUMN address varchar(200);
```

#### Drop Table

To drop an existing table:

```sh
gomig generate drop_<table_name>
```

Example:

```sh
gomig generate drop_users
```

This will generate a migration file:

```sql
-- Migration: drop_users
-- Created at: 20240101123456

DROP TABLE users;
```

#### Rename Table

To rename an existing table:

```sh
gomig generate rename_<old_table_name> <new_table_name>
```

Example:

```sh
gomig generate rename_users "customers"
```

This will generate a migration file:

```sql
-- Migration: rename_users
-- Created at: 20240101123456

ALTER TABLE users RENAME TO customers;
```

#### Drop Column

To drop columns from an existing table:

```sh
gomig generate drop_column_<table_name> <column_name> <column_name> ...
```

Example:

```sh
gomig generate drop_column_users "age" "address"
```

This will generate a migration file:

```sql
-- Migration: drop_column_users
-- Created at: 20240101123456

ALTER TABLE users
    DROP COLUMN age,
    DROP COLUMN address;
```

#### Rename Column

To rename a column in an existing table:

```sh
gomig generate rename_column_<table_name> <old_column_name>:<new_column_name>
```

Example:

```sh
gomig generate rename_column_users "email:new_email"
```

This will generate a migration file:

```sql
-- Migration: rename_column_users
-- Created at: 20240101123456

ALTER TABLE users RENAME COLUMN email TO new_email;
```

#### Modify Column

To change the type of a column in an existing table:

```sh
gomig generate modify_<table_name> <column_name:type>
```

Example:

```sh
gomig generate modify_users "age:bigint"
```

This will generate a migration file:

```sql
-- Migration: modify_users
-- Created at: 20240101123456

ALTER TABLE users
    ALTER COLUMN age TYPE bigint;
```

#### Add Index

To add an index to an existing table:

```sh
gomig generate add_index_<table_name> <index_name>:<column_name>
```

Example:

```sh
gomig generate add_index_users "users_name_idx:name"
```

This will generate a migration file:

```sql
-- Migration: add_index_users
-- Created at: 20240101123456

CREATE INDEX users_name_idx ON users (name);
```

#### Drop Index

To drop an index from an existing table:

```sh
gomig generate drop_index_<table_name> <index_name>
```

Example:

```sh
gomig generate drop_index_users "users_name_idx"
```

This will generate a migration file:

```sql
-- Migration: drop_index_users
-- Created at: 20240101123456

DROP INDEX users_name_idx;
```

### Running Migrations

To run all pending migrations:

```sh
gomig migrate --db "postgres://user:password@localhost/dbname?sslmode=disable"
```

This command will apply all pending migrations and record them in the `gomig_migrations` table.

### Library

To use GoMig as a library in your Go project, import the `migration` package and call the `Generate` function:

```go
package main

import (
    "fmt"
    "github.com/rodrigoagostin/gomig/internal/migration"
)

func main() {
    err := migration.Generate("create_users", []string{"name:varchar(200)", "email:varchar(200)", "active:boolean"})
    if err != nil {
        fmt.Println("Error generating migration:", err)
    }
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
Testing GitHub Actions
