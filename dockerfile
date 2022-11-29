# syntax=docker/dockerfile:1

FROM golang:1.17 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /gym-consistency

EXPOSE 8080

CMD [ "/gym-consistency" ]
