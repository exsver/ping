### Linters
#### Linter installation:
```sh
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin v1.15.0
go get -u github.com/Quasilyte/go-consistent
go get -u github.com/mgechev/revive
```

#### Run linters:
```sh
golangci-lint run
go-consistent -pedantic -v ./...
revive -config revive.toml -formatter friendly ./...
```