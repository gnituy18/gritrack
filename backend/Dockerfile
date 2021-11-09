FROM golang:1.16

WORKDIR /go/src/gritter
COPY ./backend .
COPY .env .

CMD ["go", "run", "main.go"]
