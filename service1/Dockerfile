FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o service service1/service1.go

CMD [ "./service" ]