FROM golang:1.11.0 as builder

WORKDIR /go/src/github.com/desertjinn/mavenlink-communicator

COPY . .

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM debian:latest

RUN mkdir /app
WORKDIR /app

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN update-ca-certificates

COPY --from=builder /go/src/github.com/desertjinn/mavenlink-communicator .

CMD ["./mavenlink-communicator"]