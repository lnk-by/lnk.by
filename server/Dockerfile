FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /bin/server ./server

EXPOSE 8080

CMD ["./server"]
