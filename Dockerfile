FROM golang:1.18

WORKDIR /app

COPY . .

RUN go mod tidy

# Install dockerize
RUN curl -sSL https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz | tar -C /usr/local/bin -xz

EXPOSE 8080

CMD ["dockerize", "-wait", "tcp://db:5432", "-timeout", "20s", "go", "run", "main.go"]
