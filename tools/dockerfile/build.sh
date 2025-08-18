#!/bin/bash
# Author: recallsong
# Email: songruiguo@qq.com

set -o errexit -o pipefail

# cd to root directory
cd $(git rev-parse --show-toplevel)

# setup base image
DOCKER_IMAGE=gohub:1.0.10

if [ -n "${DOCKER_REGISTRY}" ]; then
    DOCKER_IMAGE=${DOCKER_REGISTRY}/${DOCKER_IMAGE}
fi

# check parameters and print usage if need
usage() {
    echo "base_image.sh [ACTION]"
    echo "ACTION: "
    echo "    build       build docker image. this is default action."
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

build_image_by_buildx() {
    push=$1
    args="docker buildx build"
    args+=" -t ${DOCKER_IMAGE}"
    args+=" --label build-time='$(date '+%Y-%m-%d %T%z')' --label debian=11 --label golang=1.24"
    args+=" --platform linux/amd64,linux/arm64"
#    args+=" --progress plain"
    args+=" -f ./tools/dockerfile/Dockerfile ."
    if [[ "$push" == "true" ]]; then
        args+=" --push"
    fi

    echo "$args"
    eval "$args"
}

# push docker image
docker_login() {
    if [ -z "${DOCKER_REGISTRY}" ]; then
        echo "fail to push docker image, DOCKER_REGISTRY is empty !"
        exit 1
    fi
    if [ -n "${DOCKER_REGISTRY_USERNAME}" ]; then
        docker login -u "${DOCKER_REGISTRY_USERNAME}" -p "${DOCKER_REGISTRY_PASSWORD}" ${DOCKER_IMAGE}
    fi
}

case "$1" in
"build")
    build_image_by_buildx false
    ;;
"build-push")
    docker_login
    build_image_by_buildx true
    ;;
"image")
    echo ${DOCKER_IMAGE}
    ;;
*)
    usage
    ;;
esac
