FROM golang:1.11.0 as builder

WORKDIR /go/src/github.com/desertjinn/mavenlink-communicator

COPY . .

#RUN go get
# <- alternative way if you really don't want to use vendoring
#RUN go get -v -t -u .

# Here we're pulling in godep, which is a dependency manager tool,
# we're going to use dep instead of go get, to get around a few
# quirks in how go get works with sub-packages.
RUN go get -u github.com/golang/dep/cmd/dep
# Create a dep project, and run `ensure`, which will pull in all
# of the dependencies within this directory.
RUN dep init && dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM debian:latest

RUN mkdir /app
WORKDIR /app

RUN apt-get update
RUN apt-get install -y ca-certificates
RUN update-ca-certificates

COPY --from=builder /go/src/github.com/desertjinn/mavenlink-communicator .

CMD ["./mavenlink-communicator"]