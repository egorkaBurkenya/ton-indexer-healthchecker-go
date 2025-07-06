# Stage 1: Build the Go Healthchecker
FROM golang:1.21-alpine AS builder

# Install git to allow cloning the repository
RUN apk add --no-cache git

WORKDIR /src

# Clone the healthchecker repository
RUN git clone --depth 1 https://github.com/egorkaBurkenya/ton-indexer-healthchecker-go.git .

# Build a static binary
RUN CGO_ENABLED=0 go build -o /healthchecker .

# Stage 2: Create a minimal final image
FROM scratch
COPY --from=builder /healthchecker /healthchecker
ENTRYPOINT ["/healthchecker"]
