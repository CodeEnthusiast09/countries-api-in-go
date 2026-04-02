# =============================================================================
# STAGE 1: Build the Go binary
# =============================================================================
# We use a full Go image here because we need the Go compiler and all build
# tools. The "alpine" variant is just a lighter Linux that still has everything
# we need to compile Go code.
FROM golang:1.23-alpine AS builder

# Install git because some Go module downloads require it
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first. Docker caches each line as a "layer".
# By copying go.mod and go.sum before the rest of the code, if your code
# changes but your dependencies don't, Docker skips re-downloading them.
COPY go.mod go.sum ./
RUN go mod download

# Now copy all your source code
COPY . .

# Build the binary.
# CGO_ENABLED=0  → pure Go binary, no C library dependencies (important for Alpine)
# GOOS=linux     → build for Linux (Railway runs Linux containers)
# GOARCH=amd64   → 64-bit Intel/AMD (Railway's servers are amd64)
# -ldflags="-s -w" → strip debug info to make the binary smaller
# -o server      → name the output binary "server"
# ./cmd/         → the entry point package (where your main.go lives)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o server ./cmd/

# =============================================================================
# STAGE 2: Runtime image
# =============================================================================
# We start fresh from a clean Alpine image. We do NOT use the golang image here
# because we only need to RUN the binary, not compile it.
# This makes the final image much smaller (< 50MB vs ~300MB+).
FROM alpine:3.19

# ca-certificates → needed for your app to make HTTPS calls (to RestCountries
#                   API, exchange rate API, etc.)
# curl            → needed to download the Atlas CLI binary below
RUN apk add --no-cache ca-certificates curl

# Download the Atlas CLI binary.
# This is the exact same binary you have installed locally on your machine.
# By installing it here, atlasexec can find it at /usr/local/bin/atlas,
# AND our entrypoint script can call it directly.
RUN curl -sSf https://atlasgo.sh | sh

WORKDIR /app

# Copy the compiled Go binary from Stage 1
COPY --from=builder /app/server .

# Copy your SQL migration files into the container.
# The entrypoint script will point Atlas at this directory.
COPY migrations/ ./migrations/

# Copy the entrypoint script
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# Tell Docker (and Railway) that this app listens on port 3000.
# This does NOT open the port — it's documentation for the container runtime.
EXPOSE 3000

# This is the command that runs when the container starts.
# The entrypoint script will run migrations, THEN start the server.
ENTRYPOINT ["./entrypoint.sh"]
