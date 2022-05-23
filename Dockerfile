FROM golang:1.17 as builder

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make graph
RUN make ftp

FROM scratch

WORKDIR /app

COPY --from=builder /app/build .

EXPOSE 8080

CMD [ "/app/graphserver" ]