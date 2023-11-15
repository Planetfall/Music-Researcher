[![Go Reference](https://pkg.go.dev/badge/github.com/planetfall/music-researcher.svg)](https://pkg.go.dev/github.com/planetfall/music-researcher)
[![Go Report Card](https://goreportcard.com/badge/github.com/planetfall/music-researcher)](https://goreportcard.com/report/github.com/planetfall/music-researcher)
[![codecov](https://codecov.io/gh/planetfall/music-researcher/graph/badge.svg?token=QWPH8FP2BO)](https://codecov.io/gh/planetfall/music-researcher)
[![Tests](https://github.com/planetfall/music-researcher/actions/workflows/music-researcher.yml/badge.svg)](https://github.com/planetfall/music-researcher/actions/workflows/music-researcher.yml)
[![Release](https://img.shields.io/github/release/planetfall/music-researcher.svg?style=flat-square)](https://github.com/planetfall/music-researcher/releases)

# Music Researcher

gRpc microservice that searchs music in Spotify API

## Protobuf

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./api/music_researcher.proto

## Local development

Use local genproto
```
go mod edit -replace github.com/planetfall/genproto=../genproto
```

Use local framework
```
go mod edit -replace github.com/planetfall/framework=../framework
```

## Run

Run the server
```
go run ./cmd/server/main.go --env development
```

Run the client
```
go run ./cmd/client/main.go --host localhost:8080 --tls=false
```

## Tests

Run the tests
```
go test ./...
```

Run the tests with coverage
```
go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
```

Print the coverage in HTML
```
go tool cover -html=coverage.out
```

## Lint

Report card
```
goreportcard-cli
```

Golang Lint
```
golangci-lint run
```

## Release

```
go mod tidy
go test ./...
git commit -m "release"
git tag v0.1.0
git push origin v0.1.0
GOPROXY=proxy.golang.org go list -m github.com/planetfall/gateway@v0.1.0
```
