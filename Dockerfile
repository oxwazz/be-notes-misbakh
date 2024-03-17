FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

ENTRYPOINT [ "go", "run",  "main.go"]
EXPOSE 1323