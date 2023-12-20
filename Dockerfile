# =========================================================================== #
# Rootless containers for Go microservices.
# =========================================================================== #
# Usage:
#  SERVICE  name of service to build, found in `svc/$SERVICE`.
#
# Example:
#  docker build --build-arg SERVICE=nats-jetstream-statsd \
#    nats-jetstream-statsd:\
#    66dd84f8195d162815d29268c6f953c0fc0b65cc .
# =========================================================================== #
# STEP 1: Builder
# --------------------------------------------------------------------------- #
# Builder contains the Go compiler, build related utilities and system files
# nessecary to build a service, to be consumed by a future container.
# --------------------------------------------------------------------------- #
FROM golang:1.21 AS builder

ARG SERVICE
WORKDIR /build

# Copy go modules configuration and download early to take advantage of layer
# caching between builds and where no dependencies have changed.
COPY go.mod go.sum /build/
RUN go mod download

COPY . .

ARG REVISION=unknown

RUN CGO_ENABLED=0 go build -o /service -trimpath -ldflags="-X 'gitlab.com/laserdigital/platform/go/pkg/build.revision=${REVISION}'" svc/${SERVICE}/cmd/main.go


# --------------------------------------------------------------------------- #
# STEP 2: Result
# --------------------------------------------------------------------------- #
# Result is a rootless container containing only the service and supporting
# system files.
# --------------------------------------------------------------------------- #
FROM scratch

WORKDIR /

# run service as nobody:nobody.
USER 65534:65534
ENV UID=65534 GID=65534 USER=nobody GROUP=nobody

# establish the base set of files for the root and nobody user.
COPY --chown=0:0 build/static/etc /etc

# copy certificate authroities and timezone information from the builder.
COPY --chown=0:0 --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=0:0 --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# copy the service binary from the builder.
COPY --chown=0:0 --from=builder /service /service

ENTRYPOINT ["/service"]
