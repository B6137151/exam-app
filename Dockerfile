FROM golang:1.18

WORKDIR /app

COPY . .

RUN go mod tidy

# Add Dockerize
RUN apt-get update && apt-get install -y wget
RUN wget -O dockerize.tar.gz https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz
RUN tar -C /usr/local/bin -xzvf dockerize.tar.gz

EXPOSE 8080

CMD ["dockerize", "-wait", "tcp://db:5432", "-timeout", "20s", "go", "run", "main.go"]
