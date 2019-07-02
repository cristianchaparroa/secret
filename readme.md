# Secret Server

The following is a simple implementation of a secret server.


## Requirements

To use the secret server you should have at least:

- Golang
- git
- docker

## Installation

```
git clone github.com/cristianchaparroa/secret
cd secret
export GO111MODULE=on
go mod download
go mod tidy
go build ...
```

## Deployment

```
docker-compose up
```

## Test

```
go test ./...  -cover
```

### Project status

The following is the current state in the different branches:

Master: [![Build Status](https://travis-ci.org/cristianchaparroa/secret.svg?branch=master)](https://travis-ci.org/cristianchaparroa/secret)
