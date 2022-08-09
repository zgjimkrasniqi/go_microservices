## This will be the dokcer file that is going to tell the docker
## compose how to build an image

## base go image
#FROM golang:1.18-alpine as builder
## We are going to build our code in this image and then build
## a smaller image

## This is on the Docker image we are building
#RUN mkdir /app

## Copy everything from the current folder
#COPY . /app

## Set the working directory
#WORKDIR /app

## CGO_ENABLED=0 we are not using any C libraries, only the standard library
#RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

## Run the chmod +x (executable command) to the /app/brokerApp to make sure that is executable
#RUN chmod +x /app/brokerApp

# Build a tiny Docker Image
FROM alpine:latest

RUN mkdir /app

#COPY --from=builder /app/brokerApp /app
COPY brokerApp /app

CMD ["/app/brokerApp"]

# This should first build all of the code on one Docker image and then create a much
# smaller image and just copy over the executable
