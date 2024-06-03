FROM golang:1.22-alpine as builder
COPY . /build
WORKDIR /build

RUN go build -o out/webhook cmd/webhook/webhook.go

FROM alpine

COPY --from=builder /build/out/webhook /app/webhook

ENTRYPOINT ["/app/webhook"]