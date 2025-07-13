## Schema Design

We are using a tool called as [dbdiagram](https://dbdiagram.io/d/686eb1aff413ba35081ad5e9) to generate the schema relation and the SQL commands.

```sql
Table account as A {
    id bigserial [pk]
    owner varchar [not null]
    balance bigint [not null]
    currency varchar [not null]
    created_at timestamptz [not null, default: `now()`]

    Indexes {
        owner
    }
}

Table entries {
    id bigserial [pk]
    account_id bigint [ref: > A.id]
    amount bigint [not null]
    created_at timestamptz [not null, default: `now()`]

    Indexes {
        account_id
    }
}

Table transfers {
    id bigserial [pk]
    from_account_id bigint [ref: > A.id]
    to_account_id bigint [ref: > A.id]
    amount bigint
    created_at timestamptz [default: `now()`]

    Indexes {
        from_account_id
        to_account_id
        (from_account_id, to_account_id)
    }
}
```

## Docker setup

### Pulling an image

```bash
# Pulling an image
docker pull postgres:12-alpine
12-alpine: Pulling from library/postgres
88ec26c39e97: Pull complete 
ed2ed9ea56ac: Pull complete 
8d659638d6d6: Pull complete 
065be7c18a34: Pull complete 
f11fbb84cbf5: Pull complete 
f91df25c1740: Pull complete 
5117fda422af: Pull complete 
30292b49ff47: Pull complete 
1c931efcfd0d: Pull complete 
52f827f72350: Pull complete 
23e4ac430039: Pull complete 
Digest: sha256:7c8f4870583184ebadf7f17a6513620aac5f365a7938dc6a6911c1d5df2f481a
Status: Downloaded newer image for postgres:12-alpine
docker.io/library/postgres:12-alpine

# Listing available images
docker images
REPOSITORY   TAG         IMAGE ID       CREATED        SIZE
postgres     12-alpine   7c8f48705831   7 months ago   368MB
busybox      latest      f85340bf132a   9 months ago   6.02MB
```

### Creating a postgres container

```bash
# Pulling and creating a postgres container
# -d: running in detached mode
# -e: setting an evironment variable inside the container
# --name: name the container
docker run --name simple-banking-db -e POSTGRES_PASSWORD=password -d postgres
Unable to find image 'postgres:latest' locally
latest: Pulling from library/postgres
37259e733066: Pull complete 
183c1a68f8aa: Pull complete 
ad9bf12774e0: Pull complete 
62d6ecd3917a: Pull complete 
fd3187d6acde: Pull complete 
6c7475e9a1eb: Pull complete 
a3e06971819c: Pull complete 
8b584fe980c8: Pull complete 
7e8924348af8: Pull complete 
e0da14ecf374: Pull complete 
7683a1632945: Pull complete 
573b5990cac7: Pull complete 
47d365097430: Pull complete 
34ed8d648854: Pull complete 
Digest: sha256:3962158596daaef3682838cc8eb0e719ad1ce520f88e34596ce8d5de1b6330a1
Status: Downloaded newer image for postgres:latest
5fc3f55faed34934a014ad8ee126ffa4858a7b44c6d51a35241815fd8476b564

# Available images
docker images
REPOSITORY   TAG         IMAGE ID       CREATED        SIZE
postgres     latest      3962158596da   4 weeks ago    640MB
postgres     12-alpine   7c8f48705831   7 months ago   368MB
busybox      latest      f85340bf132a   9 months ago   6.02MB

# Running container
docker ps
CONTAINER ID   IMAGE      COMMAND                  CREATED         STATUS         PORTS      NAMES
5fc3f55faed3   postgres   "docker-entrypoint.s…"   7 seconds ago   Up 6 seconds   5432/tcp   simple-banking-

# Deleting a docker image
docker rmi <image-hash>
```

### Accessing the container via terminal

```bash
# Pulling and creating a postgres container
# -d: running the container in detached mode
# -e: setting an evironment variable inside the container
# --name: name the container
# -p: port mapping
docker run --name simple-banking-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine

# Accessing the running container
docker exec -t simple-banking-db psql -U root
# -U: specify a user

# Access docker logs
docker logs simple-banking-db
```

## DB schema migration

```shell
brew install golang-migrate
```

Follow [this](https://github.com/golang-migrate/migrate/blob/master/GETTING_STARTED.md) to get started.

Execute the command inside the project folder
```bash
migrate create -ext sql -dir db/migration -seq init_schema
```

Accessing postgres
```bash
docker exec -it simple-banking-db /bin/sh
/ # createdb --username=root --owner=root simple_bank
/ # psql simple_bank
/ # dropdb simple_bank
```

NOTE: `createbd` is a wrapper over `CREATE DATABASE`.

Add this to a Makefile
```bash
docker exec -it postgres12 createdb --username=root --owner=root simple_bank
```

Access the database console w/o going through the shell
```bash
docker exec -it postgres12 psql -U root simple_bank
```

> Use history command to lookup previously executed commands

```
ombalapure@Oms-MacBook-Air simple-bank % history | grep docker
 1175  docker ps
 1176  docker exec -it simple-banking-db psql -U root\n
 1177  docker exec -it simple-banking-db /bin/sh
 1178  docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine
 1179  docker exec -it postgres12 createdb --username=root --onwer=root simple_bank
 1180  docker exec -it postgres12 createdb --username=root --owner=root simple_bank
 1181  docker exec -it postgres12 psql -U root simple_bank
 1187  docker logs postgres12
```

Applying the migrations to our postgres database
```bash
migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
```

## Using ORM

We can use golang's builint sql module but
- Manual mapping SQL fields to variables
- Easy to make mistakes, not caught until runtime

GORM
- CRUD functions already implemented, less production code
- But learn writing queries using GORM functions
- Runs slowly on high load

SQLX
- This is a middle ground; its fast and easy to use
- Field mappings via query text and struct tags
- Failures won't occur until runtime

SQLC ** 
- Fast and easy to use
- Automatic code generation
- Catch SQL query errors before generating code
- Full support for Postgres; MySQL is experimental

### Installing SQLC

```bash
brew install sqlc
```

Use `sqlc init` to generate a `sqlc.yaml` file.

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "./db/query/"
    schema: "./db/migration/"
    gen:
      go:
        package: "db"
        out: "./db/sqlc"
        # sql_package: "pgx/v5"
```

Now add you CRUD queries in the `query/account.sql` and the execute `sqlc generate`.

This now creates `account.sql.go`, `db.go` and `models.go` under the  `sqlc` folder.

No need to write CRUD functions ourselves.

***

## Transactions

A single unit of work often made up of multiple db operations.

Eg: Transfer of 10 USD from account1 to account2:
- Create a transfer record with amount = 10
- Create an account entry for account1 with amount = -10
- Create an account entry for account2 with amount = +10
- Subtract 10 from the balance of account1
- Add 10 to the balance of account2

Why do we need it?
- Consistent unit of work, even in case of system failures
- Isolation b/w programs

How to write a transcation?
```sql
BEGIN;
...
COMMIT;
```

Special care needs to be taken when dealing with "update balance" step.

When account1 is reading existing balance to deduct the amount
The transaction tx-2 should not be able to read from account1 balance
Until tx-1 reads and updates it

This is done via
```sql
SELECT balance FROM accounts WHERE id = $1 FOR UPDATE;
```

So we add a new FOR UPDATE query; to get account details for an ID. And run `sqlc generate` again.

But using FOR UPDATE is not the correct solution; this will result in a deadlock.

To debug this issue we can pass a transaction name to `TransferTx` method via the context.

```go
// Transaction key
var txKey = struct{}{}

// Create transaction name
txName := fmt.Sprintf("tx %d", i+1)
// Create context with that transaction name
ctx := context.WithValue(context.Background(), txKey, txName)
result, err := store.TransferTx(ctx, TransferTxParams{
    FromAccountID: account1.ID,
    ToAccountID:   account2.ID,
    Amount:        amount,
})
```

Now this transaction key can be accessed in the function definition of `TransferTx`.
```go
store.execTx(ctx, func(q *Queries) error {
    var err error
    txName := ctx.Value(txKey)

    fmt.Println(txName, "create transfer")
})
```

Now we can find the actual root cause
```bash
tx 4 get account 1 for update
tx 5 get account 1 for update
tx 2 get account 1 for update
tx 3 get account 1 for update
```

**NOTE**: This happens because of a foreign key constraint: the entry or transfer table must verify that the referenced account row exists before proceeding. To perform this check, the transaction inserting into the child table tries to acquire a lock on the corresponding row in the accounts table. If another transaction is already holding a conflicting lock on that accounts row (e.g., due to an update or select...for update), this can result in a deadlock.

## Fixing Deadlock

In order to fix the deadlock use `FOR NO KEY UPDATE` on SELECT account.

When you do a SELECT ... FOR UPDATE in Postgres, you’re saying:
> "Give me an exclusive lock on this row — no one else can touch it for updates or FK checks unless they wait."

This lock blocks:
- Other SELECT ... FOR UPDATE
- Foreign key insert checks (which need to lock the parent row to ensure referential integrity)

Now enter: FOR NO KEY UPDATE — it's a weaker lock. It:
- Still prevents updates to the row's data
- But does not block foreign key checks that don’t plan to change the row’s primary/unique key


## Combining SELECT and UPDATE query

We can combine two queries to get account and update

```sql
-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;
```

> `sqlc.arg(amount)` will let us pass the amount param instead of balance which is a little misleading 

## Update GOTCHA

Using NO FOR UPDATE can still cause deadlocks

Transaction 1:
```sql
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
```

Transaction 2:
```sql
UPDATE accounts SET balance = balance - 100 WHERE id = 2;
UPDATE accounts SET balance = balance + 100 WHERE id = 1;
```

Both are transferring money between accounts:
- T1 locks id = 1, then tries to lock id = 2
- T2 locks id = 2, then tries to lock id = 1
- Deadlock! Each is holding a lock the other needs.

### The Fix

If all transactions always lock rows in ascending order (e.g., always id = 1 before id = 2), then you can never have a deadlock from order inversion.

> It could just as well be descending - doesn't matter.

Now in the code, add the following check

```go
if arg.FromAccountID < arg.ToAccountID {
    // Update account balance for account1
    result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
        Amount: -arg.Amount,
        ID:     arg.FromAccountID,
    })
    if err != nil {
        return err
    }

    // Update account balance for account2
    result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
        Amount: arg.Amount,
        ID:     arg.ToAccountID,
    })
    if err != nil {
        return err
    }
} else {
    // Update account balance for account2
    result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
        Amount: arg.Amount,
        ID:     arg.ToAccountID,
    })
    if err != nil {
        return err
    }

    // Update account balance for account1
    result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
        Amount: -arg.Amount,
        ID:     arg.FromAccountID,
    })
    if err != nil {
        return err
    }
}
```

The logic can be de-duplicated by moving it to a method.