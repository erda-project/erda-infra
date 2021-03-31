#!/bin/bash
# Author: recallsong
# Email: songruiguo@qq.com

set -o errexit -o pipefail

# check parameters and print usage if need
usage() {
    echo "protoc.sh TYPE [PATH]"
    echo "TYPE: "
    echo "    protocol     build message、grpc、http、form、register、client files."
    echo "    init         init module in current path."
    echo "    clean        clean result files of protocol building."
    echo "    message      build message only."
    echo "    usage        print usage."
    echo "PATH: "
    echo "    relative current path. default is *.proto files."
    exit 1
}
if [ -z "$1" ]; then
    usage
fi
PB_PATH=$2
if [ -z "$PB_PATH" ]; then
    PB_PATH=*.proto
fi

WORKDIR="$(pwd)"
PB_INCLUDES="${PB_INCLUDES} -I=/usr/local/include/"

# build protocol
build_protocol() {
    mkdir -p pb && mkdir -p client;
    protoc \
        -I=. ${PB_INCLUDES} \
        --go_out=./pb --go_opt=paths=source_relative \
        --go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
        --go-http_out=./pb --go-http_opt=paths=source_relative \
        --go-form_out=./pb --go-form_opt=paths=source_relative \
        --go-register_out=./pb --go-register_opt=paths=source_relative \
        --go-client_out=./client --go-client_opt=paths=source_relative \
        ${PB_PATH}
    goimports -w ./client/*.go ./pb/*.go
}

# clean result files of building
clean_result() {
    rm -rf ./client/*.go
    rm -rf ./pb/*.go
}

# build message only
build_message() {
    protoc \
        -I=. ${PB_INCLUDES} \
        --go_out=. --go_opt=paths=source_relative \
        ${PB_PATH}
    goimports -w *.go
}

# init module
init_module() {
    rm -rf ${WORKDIR}/*.go;
    HAS_GO_FILE=$(eval echo $(bash -c "find "${WORKDIR}" -maxdepth 1 -name *.go 2>/dev/null" | wc -l))
    if [ ${HAS_GO_FILE} -gt 0 ]; then
        echo "${WORKDIR} is not empty directory."
		exit 1
	fi
    protoc \
        -I=. ${PB_INCLUDES} \
        --go-provider_out=${WORKDIR} --go-provider_opt=paths=source_relative \
        ${PB_PATH}
    goimports -w *.go
}

case "$1" in
    "protocol")
        build_protocol
        ;;
    "init")
        init_module
        ;;
    "clean")
        clean_result
        ;;
    "message")
        build_message
        ;;
    *)
        usage
esac
