FROM golang:1.24.7-alpine3.22 AS deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.24.7-alpine3.22 AS builder
WORKDIR /app
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .
ENV CGO_ENABLED 0
WORKDIR /app/cmd/api
RUN go build -o api

FROM alpine:3.15
COPY --from=builder /app/cmd/api/api /app/api
WORKDIR /app
EXPOSE 8080
CMD ["./api"]
