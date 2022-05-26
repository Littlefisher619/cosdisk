FROM golang:1.17 as builder

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make graph
RUN make ftp
RUN make fuse

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch

WORKDIR /app

COPY --from=builder /app/build .
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8080

CMD [ "/app/graphserver" ]