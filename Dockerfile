FROM golang:1.18

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go get github.com/ory/dockertest/v3

EXPOSE 8080

CMD ["go", "run", "main.go"]
