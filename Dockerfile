FROM golang:1.22

RUN     mkdir /app
WORKDIR /app
ADD     . /app/
RUN     go build

CMD     ./protogofly