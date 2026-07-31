FROM golang:1.26

# Set the Current Working Directory inside the container
WORKDIR /go/src/app

# Copy go mod and sum files
COPY . .

#Expose the port
EXPOSE 8080

# Build the Go app
RUN go build -o main cmd/main.go

# Run the executable
CMD ["./main"]