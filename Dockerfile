FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags "-X main.version=1.0.0" -o vkv

CMD [ "./vkv" ]
