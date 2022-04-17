VERSION="v2.0.0"
GO_VERSION=`go version`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
GIT_COMMIT=`git rev-parse HEAD`
GIT_LATEST_TAG=`git describe --tags --abbrev=0`
BUILD_TIME=`date +"%Y-%m-%d %H:%M:%S %Z"`

LDFLAGS="-X 'main.Version=${VERSION}' -X 'main.GoVersion=${GO_VERSION}' \
-X 'main.GitBranch=${GIT_BRANCH}' -X 'main.GitCommit=${GIT_COMMIT}' \
-X 'main.GitLatestTag=${GIT_LATEST_TAG}' -X 'main.BuildTime=${BUILD_TIME}'"

all: buildProxima buildProxctl

buildProxima:
	go build -o cmd/server/proxima -ldflags ${LDFLAGS} -a  cmd/server/main.go

buildProxctl:
	go build -o cmd/client/proxctl -ldflags ${LDFLAGS} -a  cmd/client/main.go