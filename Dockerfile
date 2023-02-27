FROM golang:1.20.1-alpine3.16

RUN mkdir /goapp
ADD . /goapp

WORKDIR /goapp

RUN go build -o main .

EXPOSE 8080

CMD ["/goapp/main"]
