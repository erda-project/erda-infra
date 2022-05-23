#!/bin/bash
# Author: recallsong
# Email: songruiguo@qq.com

set -o errexit -o pipefail

# cd to root directory
cd $(git rev-parse --show-toplevel)

# setup base image
DOCKER_IMAGE=gohub:1.0.5

if [ -n "${DOCKER_REGISTRY}" ]; then
    DOCKER_IMAGE=${DOCKER_REGISTRY}/${DOCKER_IMAGE}
fi

# check parameters and print usage if need
usage() {
    echo "base_image.sh [ACTION]"
    echo "ACTION: "
    echo "    build       build docker image. this is default action."
    echo "    push        push docker image, and build image if image not exist."
    echo "    build-push  build and push docker image."
    echo "    image       show image name."
    echo "Environment Variables: "
    echo "    DOCKER_REGISTRY format like \"registry.example.org/username\" ."
    echo "    DOCKER_REGISTRY_USERNAME set username for login registry if need."
    echo "    DOCKER_REGISTRY_PASSWORD set password for login registry if need."
    exit 1
}
if [ -z "$1" ]; then
    usage
fi

# build docker image
build_image() {
    platforms="linux/amd64 linux/arm64"
    for platform in ${platforms}; do
        echo "building for $platform"
        internal_build_image $platform
    done
}

# due to the issue that Aliyun Docker Registry doesn't support multi-arch under one docker tag,
# we add targetplatform after original tag.
internal_build_image() {
    targetplatform=$1                              # linux/arm64
    dash_targetplatform=${targetplatform//\//-}    # linux-arm64
    image="${DOCKER_IMAGE}-${dash_targetplatform}" # {registry:}gohub:1.0.5-linux-arm64
    docker build \
        -t "$image" \
        --label "build-time=$(date '+%Y-%m-%d %T%z')" \
        --label "debian=11" \
        --label "golang=1.17" \
        --platform="$targetplatform" \
        --progress=plain \
        -f ./tools/dockerfile/Dockerfile .
}

# push docker image
push_image() {
    if [ -z "${DOCKER_REGISTRY}" ]; then
        echo "fail to push docker image, DOCKER_REGISTRY is empty !"
        exit 1
    fi
    IMAGE_ID="$(docker images ${DOCKER_IMAGE} -q)"
    if [ -z "${IMAGE_ID}" ]; then
        build_image
    fi
    if [ -n "${DOCKER_REGISTRY_USERNAME}" ]; then
        docker login -u "${DOCKER_REGISTRY_USERNAME}" -p "${DOCKER_REGISTRY_PASSWORD}" ${DOCKER_IMAGE}
    fi
    docker push "${DOCKER_IMAGE}"
}

# build and push
build_push_image() {
    build_image
    push_image
}

case "$1" in
"build")
    build_image
    ;;
"push")
    push_image
    ;;
"build-push")
    build_push_image
    ;;
"image")
    echo ${DOCKER_IMAGE}
    ;;
*)
    usage
    ;;
esac
