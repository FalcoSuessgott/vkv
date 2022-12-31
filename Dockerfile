# build
FROM golang:1.19.4-alpine3.16 AS builder

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o vkv

# deploy
FROM alpine:3.17
WORKDIR /build
COPY --from=builder /build/vkv /build/vkv

ENTRYPOINT ["./vkv"]