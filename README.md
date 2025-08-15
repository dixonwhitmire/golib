# golib

A collection of handy and useful "library functions" for Go Applications.

Packages include:

- [csvlib](csvlib/csv.go) provides readers (iterators) and writers for CSV files.
- [ctxlib](ctxlib/ctx.go) provides pre-configured contexts for use in applications.
- [datelib](datelib/date.go) formats time.Time values to ISO8601 formats.
- [iolib](iolib/io.go) provides access to text file contents in memory or via iterators.
- [loglib](loglib/logger.go) provides standard logger configurations.

## install
```shell
go get github.com/dixonwhitmire/golib@v0.10.0
```

## development

### dependencies

- [Go](https://go.dev/doc/install)
- [Just Command Runner](https://github.com/casey/just)

### documentation
launch the go doc web server and then browse to http://localhost:6060/pkg/github.com/dixonwhitmire/golib/

```shell
just docs
```

### execute tests

```shell
just test
```