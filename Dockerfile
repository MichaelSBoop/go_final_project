FROM golang:1.22.1

WORKDIR /go_final_project

COPY . .

RUN go mod download

ENV TODO_PORT=7540

ENV TODO_DBFILE=./

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 

RUN go build -o ./

CMD ["./go_final_project"]