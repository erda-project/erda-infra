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
