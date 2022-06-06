FROM --platform=$TARGETPLATFORM golang:1.17-bullseye as build

ENV PROJ_ROOT="/go/src/github.com/erda-project/erda-infra"
COPY . "${PROJ_ROOT}"
WORKDIR "${PROJ_ROOT}"

RUN apt-get update && \
    apt-get install -y --no-install-recommends unzip && \
    rm -fr /var/lib/apt/lists/*

ARG GOPROXY=https://goproxy.cn,direct
RUN go env -w GO111MODULE=on && go env -w GOPROXY="${GOPROXY}"
RUN go install golang.org/x/tools/cmd/goimports@latest

ENV PROTOC_VER=3.15.8
ENV PROTOC_ZIP=protoc-"${PROTOC_VER}"-linux-x86_64.zip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v"${PROTOC_VER}/${PROTOC_ZIP}" \
    && unzip -o ${PROTOC_ZIP} -d /usr/local bin/protoc \
    && unzip -o  ${PROTOC_ZIP} -d /usr/local 'include/*' \ 
    && rm -f  ${PROTOC_ZIP}

ENV PGG_VER=v1.27.1
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PGG_VER}

RUN --mount=type=cache,target=/root/.cache/go-build\
    --mount=type=cache,target=/go/pkg/mod \
    cd ./tools && \
    go install ./gohub && \
    gohub tools install --verbose --local


FROM --platform=$TARGETPLATFORM golang:1.17-bullseye

RUN go install golang.org/x/tools/cmd/goimports@latest
COPY --from=build /root/.gohub /root/.gohub
COPY --from=build /go/bin/* /go/bin/

WORKDIR /go

CMD gohub
