FROM golang:1.19-bullseye AS builder

ARG API_IMAGE_TAG

# Verify that the environment variables necessary exist
RUN [ ! -z "${API_IMAGE_TAG}" ] || { echo "API Docker Image Tag variable cannot be empty"; exit 1; }

ENV API_ENV=$API_ENV
ENV API_IMAGE_TAG=$API_IMAGE_TAG

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

ARG API_IMAGE_TAG

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
LABEL "api.image_tag"="${API_IMAGE_TAG}"

EXPOSE 80/tcp
ENTRYPOINT ["/go/bin/api"]
