FROM golang:1.17.5 AS GO_BUILD
ENV CGO_ENABLED 0
COPY . /app
WORKDIR /app/cmd/api
RUN go build -o api

FROM alpine:3.15
COPY --from=GO_BUILD /app /app
WORKDIR /app/cmd/api
EXPOSE 8080
CMD ["./api"]
