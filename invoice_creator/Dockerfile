FROM golang:1.12.9

RUN mkdir /app

ADD . /app

WORKDIR /app
# create ssh directory
# RUN mkdir ~/.ssh
# RUN touch ~/.ssh/known_hosts
# RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

# RUN git config --global url."https://mathantunes:cb3f3ca0401b9474901c85789ce75484309ace30@github.com/".insteadOf "https://github.com/"

# RUN go mod download
# RUN git clone "https://github.com/mathantunes/arex_project.git"

RUN go get    github.com/mathantunes/arex_project/filestore
RUN go get	github.com/mathantunes/arex_project/invoice
RUN go get	github.com/mathantunes/arex_project/queuer
RUN go get	github.com/mathantunes/arex_project/services
RUN go get	github.com/mathantunes/arex_project/validator
RUN go get	google.golang.org/grpc

RUN go build -o main github.com/mathantunes/arex_project/invoice_creator

CMD ["/app/main"]