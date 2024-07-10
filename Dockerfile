FROM golang:latest

COPY ./ ./
ENV GOPATH=/

RUN go mod download

RUN go build -o tx-notifier main.go
CMD ["./tx-notifier"]