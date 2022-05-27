// Hellofs implements a simple "hello world" file system.
package main

import (
	"context"
	"flag"
	"fmt"
	"runtime/pprof"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

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
	email := flag.String("email", "youremail@email.com", "Email")
	pwd := flag.String("password", "password", "Password")
	newUser := flag.String("newuser", "", "create user of username")
	mkDir := flag.Bool("mkdir", false, "create in mountpoint")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()

	handleExit(*mountpoint)
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            logrus.Fatal(err)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
	cf, err := cf.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("parse config file error: ", err)
		os.Exit(1)
	}
	service := cosdisk.New(cf)
	ctx := context.Background()
	if *newUser != "" {
		if _, err = service.UserLogin(ctx, *email, *pwd); err != nil {
			if _, err = service.UserRegister(ctx, *newUser, *email, *pwd); err != nil {
				fmt.Println("create user error: ", err)
				os.Exit(1)
			}
		}

	}

	log := logrus.New()

	if *mkDir {
		if _, err := os.Stat(*mountpoint); os.IsNotExist(err) {
			os.Mkdir(*mountpoint, 0777)
		}
	}

	// for test
	_, err = service.UserLogin(ctx, *email, *pwd)
	if err != nil {
		log.Fatal(err)
	}
	FS, err := server.NewFS(*email, *pwd, service)
	if err != nil {
		log.Fatal(err)
	}
	fuse.Debug = func(msg interface{}) {
		log.Debugf("%#v", msg)
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
	defer fuse.Unmount(*mountpoint)

	fmt.Println("Press Ctrl+C in Terminal to unmount")
	err = fs.Serve(c, FS)
	if err != nil {
		log.Fatal(err)
	}

}
