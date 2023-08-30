FROM golang:1.21.0

LABEL authors="D1m4"

WORKDIR /usr/src/app

COPY . .

#ENTRYPOINT ["top", "-b"]

RUN go mod tidy