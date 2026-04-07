FROM golang:1.25-alpine
WORKDIR /app
COPY go.mod go.sum ./        
RUN go mod download
COPY pb/ ./pb/   
COPY order/ ./order/ 
RUN go build -o main ./order
CMD ["./main"]