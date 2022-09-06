#Golang Image Base
FROM golang:alpine

#install git
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY /cmd .

RUN go build -o /app

EXPOSE 3000

CMD ["/app"]