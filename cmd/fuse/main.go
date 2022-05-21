// Hellofs implements a simple "hello world" file system.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	cf "github.com/Littlefisher619/cosdisk/config"
	server "github.com/Littlefisher619/cosdisk/fuse"
	cosdisk "github.com/Littlefisher619/cosdisk/service"
)

func handleExit(mountpoint string) {
	// cleanup and unmount on interrupt
	c1 := make(chan os.Signal)
	signal.Notify(c1, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c1
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		fuse.Unmount(mountpoint)
		os.Exit(0)
	}()
}

func main() {
	configPath := flag.String("config-file", "./config.toml", "The toml config file")
	mountpoint := flag.String("mount-point", "./test", "The mount point")
	flag.Parse()

	handleExit(*mountpoint)

	cf, err := cf.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("parse config file error: ", err)
		os.Exit(1)
	}
	service := cosdisk.New(cf)
	// for test
	_, err = service.UserRegister(context.Background(), "test", "test", "test")
	if err != nil {
		log.Fatal(err)
	}
	FS, err := server.NewFS("test", "test", service)
	if err != nil {
		log.Fatal(err)
	}

	c, err := fuse.Mount(
		*mountpoint,
		fuse.FSName("cosdisk"),
		fuse.Subtype("cosdisk"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	fmt.Println("Press Ctrl+C in Terminal to unmount")
	err = fs.Serve(c, FS)
	if err != nil {
		log.Fatal(err)
	}
}
