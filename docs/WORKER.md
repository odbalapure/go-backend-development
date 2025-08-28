## Async Workers

Async processing:
- Easiest way is to use `go routines` but this all happnes in server's memory.

Ideal way:
- Use message broker or background workers.
- Tasks are saved in both memory and persistent storage.
- HA: Redis sentinel or cluster
- Not tasks lost

Use case: send verification email
- Create a new user record
- Push a send verification email task to Redis queue
- Background worker picks up a task from the queue and process it

Install [asynq](https://github.com/hibiken/asynq)

Backed by redis:
- Client puts tasks on a queue
- Server pulls tasks off queues and starts a worker goroutine for each task
- Tasks are processed concurrently by multiple workers

Features:
- Guaranteed at least once execution of a task
- Scheduling of tasks
- Retries of failed tasks
- Automatic recovery of tasks in the event of a worker crash
- Weighted priority queues

And many others [here](https://github.com/hibiken/asynq?tab=readme-ov-file#features).

```go
go get github.com/hibiken/asynq
```

## Integrate the Redis worker

`asynq` needs redis to operate; start a redis container using:

```sh
docker run --name redis -p 6379:6379 -d redis:7-alpine
```

Check if the redis container is up and running:

```sh
2aa2eb1a6338551c50c342a8f973cecc2e4a287b2ead652b45510f05d4bc6c59
ombalapure@Oms-MacBook-Air simple-bank % docker exec -it redis redis-cli ping 
PONG
```

## Submitting and processing tasks

We can specify the queues along with their priorities while creating the `asynq` server:

```go
const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

server := asynq.NewServer(
    redisOpt,
    // Max no. processing of concurrent tasks
    // Retry delay for failed tasks
    // Error returned by handler is a failure etc.
    // Set priority of queues
    asynq.Config{
        Queues: map[string]int{
            QueueCritical: 10,
            QueueDefault:  5,
        },
    },
)
```

> We are naming the queues and their priorities. The name can be anything; its just a string.

Submitting the task to the "processor", here we can provide options like # of retries, interval between retries etc.:

```go
opts := []asynq.Option{
    // No. of retries
    asynq.MaxRetry(10),
    // Add a delay between retries
    asynq.ProcessIn(10 * time.Second),
    // Name of the queue
    // NOTE: Update the `processor` to listen to this queue
    asynq.Queue(worker.QueueCritical),
}
err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
```

## Send asynq task within a DB transaction

Create a `tx_create_user.go` file and add the following:

```go
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User `json:"user"`
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
```

> CreateUserTx is part of the `Store` interface.

The `AfterCreate` is the function submits a task to the `asynq` queue.

```go
// rcp_create_user.go
arg := db.CreateUserTxParams{
    CreateUserParams: db.CreateUserParams{
        Username:       req.GetUsername(),
        HashedPassword: hashedPassword,
        FullName:       req.GetFullName(),
        Email:          req.GetEmail(),
    },
    AfterCreate: func(user db.User) error {
        taskPayload := &worker.PayloadSendVerifyEmail{
            Username: user.Username,
        }
        opts := []asynq.Option{
            // No. of retries
            asynq.MaxRetry(10),
            // Add a delay between retries
            asynq.ProcessIn(10 * time.Second),
            // Name of the queue
            // NOTE: Update the `processor` to listen to this queue
            asynq.Queue(worker.QueueCritical),
        }
        return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
    },
}
```
