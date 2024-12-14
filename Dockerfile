FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN apt-get update && apt-get install -y docker.io

RUN go mod download

COPY server ./server

RUN go build -o /app/sync /app/server/main.go

ENV FOLDER="/app/files"

EXPOSE 5000

CMD ["/app/sync"]
