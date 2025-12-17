# Build Stage
FROM golang:1.23 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o agent ./cmd/agent

# Distroless Runtime Stage
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app/agent .
USER 65532:65532

ENTRYPOINT ["/agent"]
