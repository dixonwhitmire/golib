# docs launches a local godoc web serer
docs:
    godoc -http=:6060

# format formats golib source code using go fmt
format:
    go fmt ./...

# test executes unit tests for the project
test:
    go test ./...