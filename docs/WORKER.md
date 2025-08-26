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
