
FROM golang:alpine
# for sqlite3
RUN apk add --no-cache --virtual .build-deps \
    ca-certificates \
    gcc \
    g++
WORKDIR /go/src/mymonitor
COPY . .
RUN go env && go build ./cmd/monitoringTool/main.go


FROM alpine:latest
WORKDIR /go/src/mymonitor
COPY --from=0 /go/src/mymonitor ./
VOLUME /go/src/mymonitor/data
CMD ["./main"] 