
FROM golang:alpine
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
# for sqlite3
RUN apk add --no-cache --virtual .build-deps \
    ca-certificates \
    gcc \
    g++
WORKDIR /go/src/sqlite3_alpine
COPY . .
RUN go env && go build ./cmd/monitoringTool/main.go

FROM alpine:latest
WORKDIR /go/src/sqlite3_alpine
COPY --from=0 /go/src/sqlite3_alpine ./
# VOLUME /go/src/sqlite3_alpine/data
ENTRYPOINT ./main