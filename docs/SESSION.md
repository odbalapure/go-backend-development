## SESSION

Auth tokens are short lived, its a bad UX to make users login after every 15m. So we create refresh tokens that live for few weeks.

Create a new session table and apply a new schema

> migrate create -ext sql -dir db/migration -seq add_sessions

Apply this [new schema](../db/migration/000003_add_sessions.up.sql) using `make migrateup`.

Now generate DB CRUD functions using `make sqlc`.

### Update login handler

Add the logic to create refresh token and create session in the database.

```go
accessToken, accessPayload, err := server.tokenMaker.CreateToken(
    user.Username,
    server.config.AccessTokenDuration,
)
if err != nil {
    ctx.JSON(http.StatusInternalServerError, errorResponse(err))
    return
}

refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
    user.Username,
    server.config.RefreshTokenDuration,
)
if err != nil {
    ctx.JSON(http.StatusInternalServerError, errorResponse(err))
    return
}

// NOTE: This creates a "session record" with the "refresh token UUID"
session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
    ID:           refreshPayload.ID,
    Username:     user.Username,
    RefreshToken: refreshToken,
    UserAgent:    ctx.Request.UserAgent(),
    ClientIp:     ctx.ClientIP(),
    IsBlocked:    false,
    ExpiresAt:    refreshPayload.ExpiredAt,
})
```

> Update the controller mocks `make mockgen` and the test cases accordingly.

### Create a new handler

Create a new handler that will only accept a refresh token. Sending an "access token" to generate a "new access token" won't work because:
- We had created a session record using **refresh token ID**.
- This way only refresh tokens are allowed to get a new access tokens not.

Before sending this token we need to have few more checks:

```go
if session.IsBlocked {
    ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("blocked session")))
    return
}

if session.Username != refreshPayload.Username {
    ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("incorrect session user")))
    return
}

if time.Now().After(session.ExpiresAt) {
    ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("expired session")))
    return
}
```
