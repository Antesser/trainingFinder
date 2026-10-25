# pgxpool-transactor

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

pgxpool-transactor is a pgxpool wrapper to execute code fragments within postgres transaction.

Easiest way to get the perfect repository.

## Installation
```bash
go get github.com/golangmonster/pgxpool-transactor/v1
```

## Usage
```go
import pgxtransactor "github.com/golangmonster/pgxpool-transactor"


type repository struct {
	pool *pgxtransactor.Pool
	pgxpool.Transactor
}

func New(pool *pgxpool.Pool) *repository {
	return &repository{
		pool:       pool,
		Transactor: pool,
	}
}
```