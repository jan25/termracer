FROM golang:latest

WORKDIR /app/src
ADD . .
RUN go build -o server
ENTRYPOINT [ "./server" ]

EXPOSE 8080