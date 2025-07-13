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