FROM golang:1.19-alpine as builder

WORKDIR /backend
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -v -o server

FROM alpine:3.15
COPY --from=builder /backend/server /backend/server

CMD ["/backend/server"]
