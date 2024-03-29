#!/bin/bash

set -o errexit -o pipefail

# check parameters and print usage if need
usage() {
    echo "protoc.sh ACTION"
    echo "ACTION: "
    echo "    build        build proto to go files."
    echo "    clean        clean result files of protocol building."
    exit 1
}
if [ -z "$1" ]; then
    usage
fi

#PACKAGE_PATH=$(sed -n '/^module \(.*\)/p' go.mod)
#PACKAGE_PATH=${PACKAGE_PATH#"module "}
PACKAGE_PATH="github.com/erda-project/erda-infra/providers/component-protocol/protobuf/proto-go"

# build protocol
build_protocol() {
    BASE_PATH="../proto"
    if [ -n "${MODULE_PATH}" ]; then
        BASE_PATH="${BASE_PATH}/${MODULE_PATH}"
        echo "base path: $BASE_PATH"
        unset GEN_ALL_IMPORTS
    fi
    MODULES=$(find "${BASE_PATH}" -type d);
    for path in ${MODULES}; do
        HAS_PROTO_FILE=$(eval echo $(bash -c "find "${path}" -maxdepth 1 -name *.proto 2>/dev/null" | wc -l));
        if [ ${HAS_PROTO_FILE} -gt 0 ]; then
            if [ -z "$(echo ${path#../proto})" ]; then
                continue; # skip ../proto
            fi
            MODULE_PATH=${path#../proto/};
            echo "build module ${MODULE_PATH}";
            mkdir -p ${MODULE_PATH}/pb;
            mkdir -p ${MODULE_PATH}/client;
            gohub protoc protocol \
                 --include=../proto \
                 --msg_out="${MODULE_PATH}/pb" \
                 --service_out="${MODULE_PATH}/pb" \
                 --client_out="${MODULE_PATH}/client" \
                 --validate=true \
                 --json=true \
                 --json_opt=emit_defaults=true \
                 --json_opt=allow_unknown_fields=true \
                 ${path}/*.proto
            if [ -n "$GEN_ALL_IMPORTS" ]; then
                echo "_ \"${PACKAGE_PATH}/${MODULE_PATH}/pb\"" >> all.go
            fi
        fi;
    done;
    echo "";
    echo "build all proto successfully !";
}

# clean result files of building
clean_result() {
    MODULES=$(find "../proto" -type d);
    for path in ${MODULES}; do
        HAS_PROTO_FILE=$(eval echo $(bash -c "find "${path}" -maxdepth 1 -name *.proto 2>/dev/null" | wc -l));
        if [ ${HAS_PROTO_FILE} -gt 0 ]; then
            if [ -z "$(echo ${path#../proto})" ]; then
                continue; # skip ../proto
            fi
            MODULE_PATH=${path#../proto/};
            rm -rf ${MODULE_PATH}
        fi;
    done;
    rm -rf all.go
}

case "$1" in
    "build")
        build_protocol
        ;;
    "clean")
        clean_result
        ;;
    *)
        usage
esac
