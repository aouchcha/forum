FROM golang:1.22.3

LABEL maintainer="mohssinaynaou874@gmail.com"
LABEL version="1.0"
LABEL description="Go construite avec Docker le site est forum"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]