FROM golang:1.19-alpine

RUN apk update && \
  apk --no-cache add git tzdata bash

RUN apk add build-base

WORKDIR /app

COPY internal /app/internal
COPY cmd/pismo .
COPY go.mod go.sum ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o pismo

CMD ["/app/pismo"]