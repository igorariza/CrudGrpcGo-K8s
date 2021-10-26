#syntax=docker/dockerfile:1
FROM golang:alpine AS build
ENV GOPROXY=https://proxy.golang.org

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /api

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . .

RUN go build  ./cmd/user
EXPOSE 8080

CMD ["./user"]