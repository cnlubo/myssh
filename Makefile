BUILD_VERSION   := "0.0.1"
BUILD_DATE      := $(shell date '+%Y-%m-%d %H:%M:%S')
COMMIT_SHA1     := $(shell git rev-parse --short HEAD)

VERSION_PKG     := github.com/cnlubo/myssh
DEST_DIR        := dist
APP             := myssh

all:
	gox -osarch="darwin/amd64 linux/amd64" \
        -output='${DEST_DIR}/${APP}_{{.OS}}_{{.Arch}}' \
    	-ldflags   "-X '${VERSION_PKG}/version.Version=${BUILD_VERSION}' \
                            -X '${VERSION_PKG}/version.BuildTime=${BUILD_DATE}' \
                            -X '${VERSION_PKG}/version.GitCommit=${COMMIT_SHA1}' \
                            -w -s" \
                            ./cmd

clean:
	rm -rf ${DEST_DIR}

.PHONY : all release clean install

.EXPORT_ALL_VARIABLES:

GO111MODULE = on
