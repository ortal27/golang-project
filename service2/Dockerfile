FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o service service2/service2.go

CMD [ "./service" ]