
FROM golang:1.18


WORKDIR /app


COPY . .


RUN go mod tidy


EXPOSE 8080


CMD ["go", "run", "main.go"]
