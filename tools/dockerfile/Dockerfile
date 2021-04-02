FROM unitedwardrobe/golang-librdkafka:alpine3.12-golang1.15.0-librdkafka1.5.0 as build

RUN apk update \
&& apk add tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone \
&& apk add bash git

RUN mkdir -p "$GOPATH/src/github.com/erda-project/erda-infra/"
COPY . "$GOPATH/src/github.com/erda-project/erda-infra/"
WORKDIR "$GOPATH/src/github.com/erda-project/erda-infra/"

RUN go install ./tools/protoc-gen-go-grpc && \
    go install ./tools/protoc-gen-go-http && \
    go install ./tools/protoc-gen-go-form && \
    go install ./tools/protoc-gen-go-client && \
    go install ./tools/protoc-gen-go-provider && \
    go install ./tools/protoc-gen-go-register 

FROM unitedwardrobe/golang-librdkafka:alpine3.12-golang1.15.0-librdkafka1.5.0

RUN apk update \
&& apk add tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone \
&& apk add curl \
&& apk add bash make git 

RUN wget -O protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.15.6/protoc-3.15.6-linux-x86_64.zip && \
    mkdir protoc && \
    unzip protoc.zip -q -d protoc && \
    rm -f protoc.zip && \
    mv protoc/bin/protoc /bin/ && \
    rm -rf protoc

COPY --from=build "$GOPATH/src/github.com/erda-project/erda-infra/tools/include/" "/usr/local/include/"
COPY --from=build "$GOPATH/bin/" "/bin/"
COPY --from=build "$GOPATH/src/github.com/erda-project/erda-infra/tools/protoc.sh" "/bin/"
