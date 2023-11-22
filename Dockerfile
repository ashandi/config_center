FROM golang:1.21-alpine

WORKDIR /app

COPY ./go.mod ./
RUN go mod download

COPY . .
RUN go build -o "build.out" ./cmd

ENTRYPOINT ./build.out