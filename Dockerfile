
FROM golang:alpine
# for sqlite3
RUN apk add --no-cache --virtual .build-deps \
    ca-certificates \
    gcc \
    g++
WORKDIR /go/src/mymonitor 
# WORKDIR /app
COPY . .
RUN go env && go build ./cmd/monitoringTool/main.go
VOLUME /go/src/mymonitor/data
# copy config ... /app/data

EXPOSE 8000
CMD ["./main"] 

# FROM alpine:latest
# WORKDIR /go/src/mymonitor
# COPY --from=0 /go/src/mymonitor ./
# VOLUME /go/src/mymonitor/data
# CMD ["./main"] 