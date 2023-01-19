# syntax=docker/dockerfile:1
FROM golang:1.19-alpine

WORKDIR /app
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /main

RUN chown -R appuser:appgroup /app

USER 1000

EXPOSE 8080

CMD [ "/main" ]

