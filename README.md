# Telo

## Preparation

- Install Nodejs and yarn

- Ensure Golang already installed on your system and `GOPATH` already set.

- Install `dep`, `go get -u github.com/golang/dep/cmd/dep`

- Install `statik`, `go get -u github.com/rakyll/statik/...`

- Execute command `dep ensure`

- Execute `yarn` and `go generate` (execure go generate if you want to update statik files)

- Running application with `go install` after this run `telo -c dev.json` -> make sure $GOPATH/bin already set in environtment variable

[DEMO SITE](http://cassavahq.com)