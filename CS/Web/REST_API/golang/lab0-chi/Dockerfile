# ---- Build stage
FROM golang:1.21-alpine as builder
# This is like cd /app from inside the container
WORKDIR /app
# Copy dependencies files
COPY go.mod go.sum ./
# To install dependencies
RUN go mod download
# Aren't we copying go.mod and go.sum again with this command?
# does this command support copying directories as well?
COPY . .
# I have not idea what this command does!
# I am used to run: go build -o name ./cmd/api
# I understand different OSs need different build command
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# ---- Run stage
FROM alpine:latest
COPY --from=builder /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]
