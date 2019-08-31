FROM golang:1.12.9

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main github.com/mathantunes/arex_project/storage

CMD ["/app/main"]