# ---- Build stage
FROM golang:alpine as builder
# This is like cd /app from inside the container
WORKDIR /app
# Copy dependencies files
COPY go.mod go.sum ./
# To install dependencies
RUN go mod download
# This won't copy go.mod and go.sum again
# thanks to cache layers :)
COPY . .
# Build command
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# ---- Run stage
FROM alpine:latest
# we'll copy the static binary
COPY --from=builder /app/server /app/server
EXPOSE 3000
CMD ["/app/server"]
