FROM golang:1.24

WORKDIR /app

COPY . .
COPY ./server/GeoLite2-Country.mmdb /app/server/
COPY ./ui /app/ui
COPY ./landingpages /app/landingpages

RUN go mod download
RUN go build -o /bin/server ./server

EXPOSE 8080

CMD ["./server"]
