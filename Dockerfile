FROM golang:1.17-alpine

WORKDIR /app

COPY . .
RUN go env -w GOPROXY="https://goproxy.cn,direct"
RUN go env
RUN go mod download

RUN go build -o ftpserver cmd/ftpserver/main.go
RUN go build -o graphserver cmd/graphserver/main.go
EXPOSE 8080

CMD [ "/app/graphserver" ]