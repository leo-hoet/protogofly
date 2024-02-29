FROM golang:1.22

RUN     mkdir /app
WORKDIR /app
ADD     go.mod main.go /app/
RUN     go build

CMD     ./protofly