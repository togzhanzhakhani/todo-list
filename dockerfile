FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /todo-list

FROM alpine:latest  

COPY --from=builder /todo-list /todo-list

EXPOSE 8080

ENTRYPOINT ["/todo-list"]
