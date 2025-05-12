FROM golang:alpine AS builder

WORKDIR /build

#ADD go.mod .

COPY . .

RUN go build -o das main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/das /build/das
COPY --from=builder /build/config.yaml /build/config.yaml
COPY --from=builder /build/migrations/. /build/migrations/.

CMD [". /das"]