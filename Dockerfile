# Dockerfile References: https://docs.docker.com/engine/reference/builder/
 
# Start from the latest golang base image
FROM golang:latest as builder
 
# Add Maintainer Info
LABEL maintainer="Han <hauhancheam@pingspace.co>"
 
# Set the Current Working Directory inside the container
WORKDIR /go/src/leave-order
 
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
 
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

######## Start a new stage from scratch #######
FROM alpine:latest  
 
RUN apk --no-cache add ca-certificates
 
WORKDIR /go/src/leave-order
 
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/leave-order .
 
# Expose port 8081 to the outside world
EXPOSE 8081
 
# Command to run the executable
CMD ["./main"]