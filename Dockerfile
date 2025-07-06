# Stage 1: Build the Go Healthchecker
FROM golang:1.18-alpine AS builder

WORKDIR /src

COPY . .

# Build a static binary
RUN CGO_ENABLED=0 go build -o /healthchecker .

# Stage 2: Create a minimal final image
FROM scratch
COPY --from=builder /healthchecker /healthchecker
ENTRYPOINT ["/healthchecker"]
