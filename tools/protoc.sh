#!/bin/bash
# Author: recallsong
# Email: songruiguo@qq.com

set -o errexit -o pipefail

# check parameters and print usage if need
usage() {
    echo "protoc.sh TYPE [FIELS]"
    echo "TYPE: "
    echo "    protocol     build message、grpc、http、form、register、client files."
    echo "    init         init module in current path."
    echo "    clean        clean result files of protocol building."
    echo "    message      build message only."
    echo "    usage        print usage."
    echo "FIELS: "
    echo "    *.proto files. default is *.proto files."
    exit 1
}
if [ -z "$1" ]; then
    usage
fi
PB_FIELS=$2
if [ -z "$PB_FIELS" ]; then
    PB_FIELS=*.proto
fi
PB_DIR=$(dirname "${PB_FIELS}")

WORKDIR="$(pwd)"
PKG_PATH=$(go run github.com/erda-project/erda-infra/tools/gopkg github.com/erda-project/erda-infra)
PB_INCLUDES="${PB_INCLUDES} -I=${PKG_PATH}/tools/include/ -I=/usr/local/include/"

# build protocol
build_protocol() {
    if [ -z "$PB_OUTPUT" ]; then
        PB_OUTPUT=${PB_DIR}
    fi
    mkdir -p ${PB_OUTPUT}/pb && mkdir -p ${PB_OUTPUT}/client;
    protoc \
        -I=${PB_DIR} ${PB_INCLUDES} \
        --go_out=${PB_OUTPUT}/pb --go_opt=paths=source_relative \
        --go-grpc_out=${PB_OUTPUT}/pb --go-grpc_opt=paths=source_relative \
        --go-http_out=${PB_OUTPUT}/pb --go-http_opt=paths=source_relative \
        --go-form_out=${PB_OUTPUT}/pb --go-form_opt=paths=source_relative \
        --go-register_out=${PB_OUTPUT}/pb --go-register_opt=paths=source_relative \
        --go-client_out=${PB_OUTPUT}/client --go-client_opt=paths=source_relative \
        ${PB_FIELS}
    HAS_GO_FILE=$(eval echo $(bash -c "find "${PB_OUTPUT}/client" -maxdepth 1 -name *.go 2>/dev/null" | wc -l));
    if [ ${HAS_GO_FILE} -gt 0 ]; then
        goimports -w ${PB_OUTPUT}/client/*.go
    fi
    HAS_GO_FILE=$(eval echo $(bash -c "find "${PB_OUTPUT}/pb" -maxdepth 1 -name *.go 2>/dev/null" | wc -l));
    if [ ${HAS_GO_FILE} -gt 0 ]; then
        goimports -w ${PB_OUTPUT}/pb/*.go
    fi
}

# clean result files of building
clean_result() {
    if [ -z "$PB_OUTPUT" ]; then
        PB_OUTPUT=${PB_DIR}
    fi
    cd ${PB_OUTPUT}
    rm -rf ./client/
    rm -rf ./pb/
}

# build message only
build_message() {
    if [ -z "$PB_OUTPUT" ]; then
        PB_OUTPUT=${PB_DIR}
    fi
    protoc \
        -I=${PB_DIR} ${PB_INCLUDES} \
        --go_out=${PB_OUTPUT} --go_opt=paths=source_relative \
        ${PB_FIELS}
    goimports -w ${PB_OUTPUT}/*.go
}

# init module
init_module() {
    HAS_GO_FILE=$(eval echo $(bash -c "find "${WORKDIR}" -maxdepth 1 -name *.go 2>/dev/null" | wc -l))
    if [ ${HAS_GO_FILE} -gt 0 ]; then
        echo "${WORKDIR} is not empty directory."
		exit 1
	fi
    protoc \
        -I=${PB_DIR} ${PB_INCLUDES} \
        --go-provider_out=${WORKDIR} --go-provider_opt=paths=source_relative \
        ${PB_FIELS}
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
