FROM golang:1.22 as builder

RUN     mkdir /app
WORKDIR /app
ADD     . /app/
RUN     go build

FROM debian:bullseye-slim
COPY --from=builder /app/protogofly .
CMD  ./protogofly