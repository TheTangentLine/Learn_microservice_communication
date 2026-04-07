FROM golang:1.25-alpine
WORKDIR /app
COPY go.mod go.sum ./        
RUN go mod download
COPY pb/ ./pb/   
COPY inventory/ ./inventory/ 
RUN go build -o main ./inventory
CMD ["./main"]