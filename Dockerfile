FROM golang:1.19-bullseye AS builder

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 65532 \
    small-user

WORKDIR /go/src/app

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ ./cmd/...

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/src/app/bin/cmd /go/bin/api

USER small-user:small-user

MAINTAINER "Conor Mc Govern <conor(at)mcgov(dot)ie>"
MAINTAINER "University Of Galway Computer Society <compsoc(at)socs(dot)nuigalway(dot)ie>"
LABEL "traefik.default.protocol"="http"
LABEL "traefik.port"="80"
LABEL "traefik.enable"="true"

EXPOSE 80/tcp
ENTRYPOINT ["/go/bin/api"]
