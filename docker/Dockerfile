FROM golang:1.17-alpine as builder

WORKDIR /app
COPY . .
RUN go mod download
RUN cd cmd/serve \
    && go build .

FROM alpine:3.15

RUN mkdir -p /opt/server/bin 

COPY --from=builder /app/cmd/serve/serve /opt/server/bin/serve

WORKDIR /opt/server

