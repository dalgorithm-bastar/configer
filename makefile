VERSION="v3.1.0"
GO_VERSION=`go version`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
GIT_COMMIT=`git rev-parse HEAD`
GIT_LATEST_TAG=`git describe --tags --abbrev=0`
BUILD_TIME=`date +"%Y-%m-%d %H:%M:%S %Z"`

LDFLAGS="-X 'main.Version=${VERSION}' -X 'main.GoVersion=${GO_VERSION}' \
-X 'main.GitBranch=${GIT_BRANCH}' -X 'main.GitCommit=${GIT_COMMIT}' \
-X 'main.GitLatestTag=${GIT_LATEST_TAG}' -X 'main.BuildTime=${BUILD_TIME}'"

all: buildProxima buildProxctl

testgo:
	go test ./... -coverprofile=cover.out
	go tool cover -func cover.out

buildProxima:
	go build -o cmd/server/proxima -ldflags ${LDFLAGS} -a  cmd/server/main.go

buildProxctl:
	go build -o cmd/client/proxctl -ldflags ${LDFLAGS} -a  cmd/client/main.go

packageAndSent:
	mkdir -p configcenter/client configcenter/server
	cp -r cmd/server/config configcenter/client
	cp -r cmd/server/config configcenter/server
	go build -o configcenter/client/proxctl -ldflags ${LDFLAGS} -a  cmd/client/main.go
	go build -o configcenter/server/proxima -ldflags ${LDFLAGS} -a  cmd/server/main.go
	tar -zcf xarch-4.0.0-rhel8.3-configcenter.tar.gz configcenter
	rm -rf configcenter