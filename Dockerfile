FROM golang:1.20.13-alpine3.19
WORKDIR /app 
COPY . .
RUN go build -o main main.go
CMD ["/app/main"]