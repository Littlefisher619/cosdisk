all: ftp graph fuse
.PHONY: ftp graph fuse clean fmt test all

ftp:
	go build -o build/ftpserver cmd/ftpserver/main.go

fuse:
	go build -o build/fuse cmd/fuse/main.go

graph:
	go build -o build/graphserver cmd/graphserver/main.go

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	rm -rf build
