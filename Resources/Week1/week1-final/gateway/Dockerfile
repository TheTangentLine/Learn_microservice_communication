FROM golang:1.25-alpine
WORKDIR /app
COPY go.mod go.sum ./        
RUN go mod download
COPY gateway/ ./gateway/ 
RUN go build -o main ./gateway
CMD ["./main"]