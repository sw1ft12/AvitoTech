FROM golang:1.22

WORKDIR /app

COPY . .

RUN go build -o app .
EXPOSE 8080

CMD ["./app"]