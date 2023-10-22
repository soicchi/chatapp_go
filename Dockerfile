FROM golang:1.21.0-bullseye as dev
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download
COPY . .
EXPOSE 3000
CMD ["go", "run", "./cmd/chatapp/main.go"]
