## Testing in Go

### Setting up the test

Go has a convention of having test files in the same directoy as the srouce code.

`main_test.go` is the entry point of the tests.

```go
package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "posgtres"
	dbSource = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

// Global variable for our tests
var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

    // Using `New()` from sqlc
	testQueries = New(conn)

	os.Exit(m.Run())
}
```

You will get the following error after runnng the test
> cannot connect to db: sql: unknown driver "posgtres"

So we need to install a postgres driver
> go get github.com/lib/pq

Keep the import to avoid the error 

```go
import (
    _ "github.com/lib/pq"
)
```

### Writing the first test

We can install an external library for assertion eg: `testify`. Read more about it, over [here](https://github.com/stretchr/testify).

Install it using
> go get github.com/stretchr/testify

Run tests where ever present and show coverage

```go
package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "tom",
		Balance:  100,
		Currency: "USD",
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
```

```bash
ombalapure@Oms-MacBook-Air simple-bank % go test -v -cover ./...
        simple-bank             === RUN   TestCreatAccount
--- PASS: TestCreatAccount (0.08s)
PASS
coverage: 6.5% of statements
ok      simple-bank/db/sqlc     0.701s  coverage: 6.5% of statements
        simple-bank/util                coverage: 0.0% of statements
```

## Mocking DB

Why need to mock the database?
- Independent tests: isolates test data to avoid conflicts
- Faster tests: reduce talking to a database
- 100% coverage: Easily setup edge cases or unexpected error like connection error

> Is it good enough to test the API with mock DB. Yes, if the real DB has integration tests.

How to mock?
- Use fake DB; store data in-memory
- Use DB stubs; Go Mock that generate and build stubs that returns hard-coded values

Install a golang mock
```go
go get github.com/golang/mock/mockgen@v1.6.0
```

The mockgen file will be present under ~/go/bin

```sh
ombalapure@Oms-MacBook-Air simple-bank % ls -l ~/go/bin 
total 122280
-rwxr-xr-x@ 1 ombalapure  staff  37816306 Jul  6 14:40 gopls
-rwxr-xr-x@ 1 ombalapure  staff  10091442 Jul 15 09:17 mockgen
```

## Creating Store interface

Since the `NewServer` takes db store object connects to the real DB. We need to replace that object with an interface.

We can create an interface ourselves but its time consuming, we can get one from `sqlc` by adding `emit_interface: true` in the sqlc.yaml file.

This generates a querier.go file with Querier interface.

```go
type Querier interface {
	AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	DeleteAccount(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	GetEntry(ctx context.Context, id int64) (Entry, error)
	GetTransfer(ctx context.Context, id int64) (Transfer, error)
	ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error)
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
}
```

Now we can make Store as interface.

```go
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}
```

Interfaces are already references so we cannot use `*Store` as the return type.

```go
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}
```

## Writing mock tests

Create a `mock` folder under the `db` folder.

```sh
ombalapure@Oms-MacBook-Air simple-bank % mockgen -help
mockgen has two modes of operation: source and reflect.

Source mode generates mock interfaces from a source file.
It is enabled by using the -source flag. Other flags that
may be useful in this mode are -imports and -aux_files.
Example:
        mockgen -source=foo.go [other options]

Reflect mode generates mock interfaces by building a program
that uses reflection to understand interfaces. It is enabled
by passing two non-flag arguments: an import path, and a
comma-separated list of symbols.
Example:
        mockgen database/sql/driver Conn,Driver
```

Create the Store mock using

```sh
mockgen -package mockdb -destination db/mock/store.go simple-bank/db/sqlc Store
```

- `-package` sets the name of the package as mockdb
- `-destination` specifies where the mock need to be created
- `simple-bank/db/sqlc` tells where the Store interace is
- `Store` is the name of the interface itself

> This will generate a `store.go` under `db/mock`.

### Adding coverage for get account controller

```go
func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	// Create a mock controller to manage mock expectations
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock database store
	store := mockdb.NewMockStore(ctrl)
	store.
		EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// Create server instance with mock store dependency
	server := NewServer(store)

	// Create response recorder to capture HTTP response
	recorder := httptest.NewRecorder()

	// Build request URL with account ID
	url := fmt.Sprintf("/accounts/%d", account.ID)
	// Create HTTP GET request
	request, err := http.NewRequest(http.MethodGet, url, nil)

	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	// Verify response body matches expected account data
	requireBodyMatchAccount(t, recorder.Body, account)
}
```

We can create a test table to cover all possible scenarios in a single function. Please refer [account_test.go file](../api/account_test.go).

NOTE: Gin runs the tests with **Debug** mode; we need configure it with to run with **Test** mode.

## Fixing test case after V2 DB migration

### Fixing integration tests

The integration/db tests will fail because the foreign key constraints are not met.
```sh
pq: insert or update on table "accounts" violates foreign key constraint "accounts_owner_fkey"
```

> So we first create a user and then an account in our tests to fix this test error

### Fixing controller/handler tests

Run the mockgen command to implement the user methods

```sh
mockgen -package mockdb -destination db/mock/store.go simple-bank/db/sqlc Store
```

## Adding config object

Since the [NewServer](../api/server.go) function now takes in a config object, we need to update the NewServer being used in the tests.

```go
func NewServer(config util.Config, store db.Store) (*Server, error) {
	// ...
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	// ...
}
```

Create a new function, as a wrapper over the NewServer.

```go
func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
```
