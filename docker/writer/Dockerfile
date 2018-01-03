FROM golang:1.9.2-alpine

WORKDIR /go/src/rightdog
COPY . .

WORKDIR cmd/writer

RUN go-wrapper install 

CMD ["go-wrapper", "run"]
