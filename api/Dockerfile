FROM golang:alpine as builder

COPY . .

RUN go build -o main main.go

CMD ["./main"]
