FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go build -o main .

CMD ["/app/main"]
