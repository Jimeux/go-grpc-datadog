FROM golang:1.19.3-bullseye AS builder
WORKDIR /go/src/github.com/Jimeux/go-grpc-datadog/svc/$SERVICE
ARG SERVICE

COPY ./proto/go ../../proto/go
COPY ./svc/$SERVICE/go.* ./
RUN go mod download

COPY ./svc/$SERVICE .
# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/service

FROM scratch

# Run as non-privileged user
COPY --from=builder /etc/passwd /etc/passwd
USER nobody

# Copy application binary
COPY --from=builder /go/bin/service /service

ENTRYPOINT ["/service"]
