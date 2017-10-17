PKG=github.com/mmcloughlin/pearl
CMD=${PKG}/cmd/pearl
GITSHAFULL=`git rev-parse HEAD`
GITSHA=`git rev-parse --short HEAD`
LDFLAGS="-X ${PKG}/meta.GitSHAFull=${GITSHAFULL} -X ${PKG}/meta.GitSHA=${GITSHA}"

install:
	go install -ldflags ${LDFLAGS} ${CMD}
