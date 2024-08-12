
# How to contribute

## Get packages
`go get github.com/takecontrolsoft/go_multi_log@v1.0.1`

## Build go exe
* to local folder - `go build -v ./...`

* to bin folder `go build -o bin/`

# How to run exe
## Open exe help
`bin/file_rename.exe /help`

# Start documentation
## To build documentation
```bash
go get golang.org/x/tools/cmd/godoc
export GOPATH=$HOME/go 
export GOROOT=/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:$GOROOT/bin
mkdir -p $GOPATH $GOPATH/src $GOPATH/pkg $GOPATH/bin
go install golang.org/x/tools/cmd/godoc@latest
godoc -http=:8081 -index
```
## Brows documentation
 http://localhost:8081/pkg/

