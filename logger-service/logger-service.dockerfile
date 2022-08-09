# Build a tiny Docker Image
FROM alpine:latest

RUN mkdir /app

COPY loggerApp /app

CMD ["/app/loggerApp"]
