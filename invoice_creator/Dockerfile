FROM golang:1.12.9

RUN mkdir /app

ADD . /app

WORKDIR /app
RUN echo "[url \"git@github.com:\"]\n\tinsteadOf = https://github.com/" >> /root/.gitconfig
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config
RUN go mod download

RUN go build -o main .

CMD ["/app/main"]