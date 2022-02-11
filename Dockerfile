FROM golang:1.17.7-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags "-X main.version=1.0.0" -o vkv

CMD [ "./vkv" ]
