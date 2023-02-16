FROM golang:1.17-alpine

WORKDIR /

RUN apk add --update yourPackageName

RUN apt update && apt install -y install build-essential && rm -rf /var/lib/apt/lists/*

COPY go.mod ./

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build ./cmd/monitoringTool/main.go

CMD [ "/gomonitor" ]