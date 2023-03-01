FROM golang:1.19.2 AS builder

WORKDIR /app

COPY . /app/

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o main

FROM scratch
WORKDIR /internal
COPY --from=builder /app/ /internal

EXPOSE 50000

ENTRYPOINT ["./main"]