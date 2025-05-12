FROM golang:alpine AS builder

WORKDIR /build

#ADD go.mod .

COPY . .

RUN go build -o das main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/das /build/das

CMD [". /das"]