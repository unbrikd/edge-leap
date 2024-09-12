FROM golang:1.22 AS builder

ARG GO_PKG
ARG APPLICATION_BUILDID
ARG APPLICATION_VERSION

WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X ${GO_PKG}/version.Version=${APPLICATION_VERSION} -X ${GO_PKG}/version.Revision=${APPLICATION_BUILDID}" -o /usr/bin/elcli

FROM alpine:3.9 AS runner
COPY --from=builder /usr/bin/elcli /usr/bin/elcli
RUN chmod +x /usr/bin/elcli

# Set the entrypoint to receive the arguments from the command line when executing the container
ENTRYPOINT ["/usr/bin/elcli"]
