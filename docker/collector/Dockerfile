FROM golang:1.9.2-alpine

WORKDIR /go/src/rightdog
COPY . .

WORKDIR cmd/collector

RUN go-wrapper install 

CMD ["go-wrapper", "run"]
