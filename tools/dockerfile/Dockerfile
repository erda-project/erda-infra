FROM unitedwardrobe/golang-librdkafka:alpine3.12-golang1.15.0-librdkafka1.5.0 as build

RUN apk update \
&& apk add tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone \
&& apk add bash git

RUN mkdir -p "$GOPATH/src/github.com/erda-project/erda-infra/"
COPY . "$GOPATH/src/github.com/erda-project/erda-infra/"
WORKDIR "$GOPATH/src/github.com/erda-project/erda-infra/"

RUN cd ./tools && \
    go install ./protoc/protoc-gen-go-grpc && \
    go install ./protoc/protoc-gen-go-http && \
    go install ./protoc/protoc-gen-go-form && \
    go install ./protoc/protoc-gen-go-client && \
    go install ./protoc/protoc-gen-go-provider && \
    go install ./protoc/protoc-gen-go-register && \
    go install ./gohub

RUN mkdir -p /root/.gohub/ && \
    cd /root/.gohub/ && \
    wget -O protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.15.6/protoc-3.15.6-linux-x86_64.zip && \
    mkdir protoc.tmp && \
    unzip protoc.zip -q -d protoc.tmp && \
    rm -f protoc.zip

FROM unitedwardrobe/golang-librdkafka:alpine3.12-golang1.15.0-librdkafka1.5.0

RUN apk update \
&& apk add tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Asia/Shanghai" > /etc/timezone \
&& apk add curl \
&& apk add bash make git 

COPY --from=build "$GOPATH/src/github.com/erda-project/erda-infra/tools/protoc/include/" "/usr/local/include/"
COPY --from=build "$GOPATH/bin/protoc-*" "/root/.gohub/"
COPY --from=build "$GOPATH/bin/gohub" "/bin/gohub"
COPY --from=build "/root/.gohub/protoc.tmp/bin/protoc" "/root/.gohub/protoc"
