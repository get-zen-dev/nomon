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
RUN go build -o main ./cmd/monitoringTool/main.go

FROM alpine:3.14

# Set the working directory to /app
WORKDIR /app

# Copy only the necessary files from the build stage
COPY --from=build /app .

# Expose port 8000
EXPOSE 8000

# Run the application
CMD ["./main"]