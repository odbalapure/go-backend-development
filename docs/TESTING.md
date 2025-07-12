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

So we need to install a postgres drive
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

> Running the main test suite also shows you the code coverage.
