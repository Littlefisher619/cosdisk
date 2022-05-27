all: ftp graph fuse
.PHONY: ftp graph fuse clean fmt test all

ftp:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/ftpserver cmd/ftpserver/main.go

fuse:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/fuse cmd/fuse/main.go

graph:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/graphserver cmd/graphserver/main.go

fmt:
	go fmt ./...

test:
	go test ./...

FUSE_MOUNTPOINT ?= /tmp/fuse_test
CONFIG_FILE ?= config.toml
benchmark:
	: '$(FUSE_MOUNTPOINT)'
	: '$(CONFIG_FILE)'
	make fuse
	./build/fuse --cpuprofile 1.prof --mkdir --newuser test --email test@email.com --password password --config-file $(CONFIG_FILE) --mount-point $(FUSE_MOUNTPOINT) 1>fuse.log 2>&1 &
	sleep 3
	FUSE_MOUNTPOINT=$(FUSE_MOUNTPOINT) go test -benchmem -run="^$$" -bench="^BenchmarkFuse$$" github.com/Littlefisher619/cosdisk/fuse
	umount $(FUSE_MOUNTPOINT)
	rmdir $(FUSE_MOUNTPOINT)

clean:
	rm -rf build
