
FROM golang:alpine
# for sqlite3
RUN apk add --no-cache --virtual .build-deps \
    ca-certificates \
    gcc \
    g++
WORKDIR /app
COPY . .
RUN go env && go build ./cmd/monitoringTool/main.go
VOLUME /app/data
# copy config ... /app/data

EXPOSE 8000

# FROM alpine:latest
# WORKDIR /app
# COPY --from=0 /app/data/ /app/index.html /app/main/ /app/icon.png /app/CloudronManifest.json ./
# VOLUME /app/data
RUN ls -al
CMD ["./main"] 

FROM golang:alpine AS build

# Install any dependencies required by your project
RUN apk add --no-cache git && \
    apk add --no-cache --virtual .build-deps \
    ca-certificates \
    gcc \
    g++

# Set the working directory to /app
WORKDIR /app

# Copy the source code of your project into the container
COPY . .

# Build the project
RUN go build ./cmd/monitoringTool/main.go

FROM alpine:3.14

# Set the working directory to /app
WORKDIR /app

# Copy only the necessary files from the build stage
COPY --from=build /app .

# Create a named volume for the /app/data directory
VOLUME ["/app/data"]

# Expose port 8000
EXPOSE 8000

# Run the application
CMD ["./main"]