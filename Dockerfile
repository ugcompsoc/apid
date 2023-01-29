FROM golang:1.19-bullseye AS builder

WORKDIR /go/src/app
COPY go.* ./
RUN go mod download

COPY cmd/ ./cmd/
#COPY internal/ ./internal/     # uncomment when development starts
RUN mkdir bin/ && go build -o bin/ ./cmd/...

FROM debian:buster-slim

RUN apt update && apt-get install -y ca-certificates && update-ca-certificates
COPY --from=builder /go/src/app/bin/cmd /go/bin/api

EXPOSE 80/tcp
ENTRYPOINT ["/go/bin/api"]