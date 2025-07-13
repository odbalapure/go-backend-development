## CI using Github Actions

Creating a job to run tests

```yml
# This workflow will build a golang project
# For more information see: https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go

name: ci-test

on:
  push:
    branches: [ "master" ]
  pull_request:
    branches: [ "master" ]

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:

    - name: Setup go 1.x
      uses: actions/setup-go@v4
      with:
        go-version: ^1.24
      id: go

    - name: Checkout code into the Go module directory
      uses: actions/checkout@v4

    - name: Test
      run: make test
```

This job will fail as postgres is not setup.

```yml
services:
    # Label used to access the service container
    postgres:
    # Docker Hub image
    image: postgres
    # Provide the password for postgres
    env:
        POSTGRES_PASSWORD: postgres
        POSTGRES_USER: root
        POSTGRES_DB: simple_bank
    # Set health checks to wait until postgres has started
    options: >-
        --health-cmd pg_isready
        --health-interval 10s
        --health-timeout 5s
        --health-retries 5
```

This workflow will fail again because migrate will not work.

Go to the [releases section](https://github.com/golang-migrate/migrate/releases) and copy the link address for [linux](https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz).

Install the golang-migrate before running `migrateup` step.

```yml
- name: Install golang-migrate
    run: |
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz
    sudo mv migrate.linux-amd64.tar.gz /usr/bin/migrate
    which migrate
```

> This movies and renames the binary file to `migrate`.
