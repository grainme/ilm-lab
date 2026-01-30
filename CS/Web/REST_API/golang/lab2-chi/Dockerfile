FROM golang:alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

FROM alpine:latest
COPY --from=builder /app/server /app/server
EXPOSE 3000
CMD ["/app/server"]
