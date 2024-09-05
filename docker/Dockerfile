FROM golang:1.22 AS builder

ARG GOOS
ARG GOARCH

WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-w -s" -o /usr/bin/elcli

FROM alpine:3.9 AS runner
COPY --from=builder /usr/bin/elcli /usr/bin/elcli
RUN chmod +x /usr/bin/elcli

# Set the entrypoint to receive the arguments from the command line when executing the container
ENTRYPOINT ["/usr/bin/elcli"]
