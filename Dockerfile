FROM golang:1.20 as builder
ENV CGO_ENABLED 0

COPY . /service

# Build the binary.
WORKDIR /service
RUN go build -o service main.go

FROM alpine:3.18
COPY --from=builder /service/service /service/service

WORKDIR /service
CMD ["./service"]

