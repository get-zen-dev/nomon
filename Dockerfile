FROM golang:latest-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . ./app

RUN go build -o ./app/cmd/monitoringTool/main.go

CMD [ "/StartMonitor" ]