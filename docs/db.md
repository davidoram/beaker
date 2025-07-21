# Database layer

## Database selection

I choose Postgres because its well known, well documented, runs pretty much anywhere and is very reliable. See [here](https://dev.to/shayy/postgres-is-too-good-and-why-thats-actually-a-problem-4imc) for more detailed discussion on this.  Needless to say as a developers you have plenty of challenges with your tech stack, so choose a database that won't give you any issues. I'm not experienced enough in other databases to recommend them, but I'm sure there are quite a few other RDBMS options that will work just as well.  

## Database creation

To support development ans unit tests we run two databases, namely `beaker_development` and `beaker_test`

The [Makefile `create-db`](../Makefile) target creates the database.  This is where you set global properties against the database.  A couple of examples that we set here are:

- `WITH OWNER postgres`: Means that the Postgres "user" called `postgres` **owns** this database, with full rights to read/write and manage the database.  In a real production system you might have other users with different rights, eg: Developers might be issued with a read/only postgres user to query the database, and by doing that we know they can't accidently change or delete the data.
- `ENCODING 'UTF8'`: This allows a broad range of languages and unicode chars to be stored.
- `LC_COLLATE='en_US.UTF-8'`: Sets the rules on how text is sorted and compared. Matching is **case sensitive** so 'a' is considered different to 'A'. We have to remember this when we run queries to find data.
- `LC_CTYPE='en_US.UTF-8'`: This sets the rules for recognizing character types (like letters, numbers, etc.).


Database migrations are a tool used to control DDL or database definition changes using version controlled scripts.  This means our data structure changes are controlled in the same way as application source code changes.

We use a tool called [`sql-migrate`](https://github.com/rubenv/sql-migrate) to do this. You can integrate this tool in many ways, for example you can run it standalone, or embed it as a library.  My preference is to run it standalone because it provides the following advantages:

- Simplifies the application code because the app can run assuming that the database is fully created, and has all the tables, indexes and other artifacts needed.
- Having a standalone app to hand on your production servers gives you extra capabilities you might need one day, for example the ability to "rollback" changes - see [`cli` options](https://github.com/rubenv/sql-migrate?tab=readme-ov-file#as-a-standalone-tool).   

## Database table design

The database stores our application data.  Programnming languages, and systems come and go, but databases have much longer lives than the applications that use them, so we need to spend extra time and care to ensure we have modelled the solution well in our database.

Its my prefernce to constrain the data as much as possible in the database, which means in our case matching data types on tables, constraining values to what we know to be valid in the domain.

In our case we have a single [database migration](db-migrations/20250716085349-create-tables.sql), which will:
- Create the `inventory` table, with columns `product_sku` and `stock_level`
- Constrain `product_sku` to lowercase letters, hyphen, underscore and numbers
- Constrain `stock_level` to be > 0

## Database environments

Run `make recreate-db` to:
- kill any connections to the database
- drop the database
- create the database
- apply all the the database migrations

We will end up with a `beaker_development` database.

Run `DB_ENV=test make recreate-db` to perform the same steps but on a `beaker_test` database.

The `beaker_development` is for development, and the `beaker_test` is for unit tests.

There is one school of opinion that running unit tests against a real database might slow your tests down, but I feel that tradeoff is worthwhile because you end up with tests that exercise real production code through all layers of your application.

## Query layer and Postgres driver selection

The `go` programming language provides a standard interface to access sql database through the [`database/sql`](https://pkg.go.dev/database/sql) package.  Its then up to the developer to choose a [Database drivers](https://github.com/avelino/awesome-go?tab=readme-ov-file#database-drivers) that will suit your application.

For out project we will use [pgx](https://github.com/jackc/pgx).

`pgx` is well supported by the `sqlc` tool that we use to generate our query bolierplate code, and can be integrated with our [Open Telemetry](./otel.md) using the [sqlc-pgx-monitoring](https://github.com/amirsalarsafaei/sqlc-pgx-monitoring). This is important at runtime, because it will allow us to see inside the running application and examine how our queries are performing. 

[`sqlc`](https://sqlc.dev) turns our SQL queries into typesafe `go` code, which our application will in turn use to access the database.  Why do we use `sqlc` instead of writing code by hand?
- It writes code better than we can?
- It performs typesafe conversions