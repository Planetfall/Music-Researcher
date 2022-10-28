FROM golang:1.18-alpine

ENV APP_ENV=production

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
RUN apk add git

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/server

RUN echo $APP_ENV

CMD app --env ${APP_ENV}
