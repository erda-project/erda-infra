ARG GO_VERSION=1.22
ARG BASE_IMAGE="golang:${GO_VERSION}-bullseye"

FROM --platform=$TARGETPLATFORM ${BASE_IMAGE} as build

ENV PROJ_ROOT="/go/src/github.com/erda-project/erda-infra"
WORKDIR "${PROJ_ROOT}"

ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on \
    CGO_ENABLED=0

# Copy go.mod and go.sum to leverage Docker cache
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy all project files
COPY . .

# Install goimports and gohub
RUN go install golang.org/x/tools/cmd/goimports@latest

# Build gohub with cache mechanism
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    cd ./tools && \
    go install ./gohub && \
    gohub tools install --verbose --local

FROM --platform=$TARGETPLATFORM ${BASE_IMAGE}

RUN go install golang.org/x/tools/cmd/goimports@latest

# Copy built tools from build stage
COPY --from=build /root/.gohub /root/.gohub
COPY --from=build /go/bin/* /go/bin/

WORKDIR /go

CMD ["gohub"]
