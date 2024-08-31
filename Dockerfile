FROM golang:1.22.1 AS builder

RUN mkdir app

WORKDIR /app

COPY . .

RUN go mod download

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o ./app .

FROM alpine

RUN adduser -S -D -H -h /app user

USER user

WORKDIR /app

COPY /web ./web

COPY *.sql .env ./

COPY --from=builder app/app ./

CMD ["./app"]