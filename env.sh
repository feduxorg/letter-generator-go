: ${GOOS:=linux}
: ${GOARCH:=amd64}
: ${GO_TAGS:=static}
: ${LIBGIT2_TAG:=release-1.5}

: ${COMMIT_HASH:=$(git rev-parse --short HEAD 2>/dev/null)}
: ${BUILD_DATE:=$(date +%FT%T%z)}
: ${VERSION:=$(git tag | sort | tail -n 1)}
: ${VERSION:=${VERSION/v/}}
