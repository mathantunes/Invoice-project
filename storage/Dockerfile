FROM golang:1.12.9

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN git config --global url."https://71d1b6c302f3aca17943ad317d552b381f455544:@github.com/".insteadOf "https://github.com/"

RUN go build -o main github.com/mathantunes/arex_project/storage

CMD ["/app/main"]