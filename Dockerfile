# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN apk add build-base

COPY *.go ./

COPY /data ./data
COPY /writeToDatabase ./writeToDatabase
COPY .gitignore ./.gitignore
COPY Dockerfile ./Dockerfile

RUN go build -o /finance

EXPOSE 7000

CMD [ "/finance" ]
