package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Littlefisher619/cosdisk/config"
	server "github.com/Littlefisher619/cosdisk/ftp"
	cosdisk "github.com/Littlefisher619/cosdisk/service"
	lib "github.com/fclairamb/ftpserverlib"
	"github.com/sirupsen/logrus"
)

var (
	ftpServer *lib.FtpServer
	driver    *server.Server
)

func main() {
	// Arguments vars
	var confFile string
	var onlyConf bool

	// Parsing arguments
	flag.StringVar(&confFile, "conf", "", "Configuration file")
	flag.BoolVar(&onlyConf, "conf-only", false, "Only create the conf")
	configPath := flag.String("config-file", "./config.toml", "The toml config file")
	flag.Parse()

	// Setting up the logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	conf := server.Config{
		Version:                  1,
		PassiveTransferPortRange: &server.PortRange{Start: 2121, End: 2131},
	}

	c, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("parse config file error: ", err)
		os.Exit(1)
	}
	service := cosdisk.New(c)
	// for test
	service.UserRegister(context.Background(), "hhh", "hhh", "hhh")
	// Loading the driver
	var errNewServer error
	driver, errNewServer = server.NewServer(&conf, service, logrus.WithField("component", "driver"))

	if errNewServer != nil {
		logger.Error("Could not load the driver", "err", errNewServer)

		return
	}

	// Instantiating the server by passing our driver implementation
	ftpServer = lib.NewFtpServer(driver)

	// Overriding the server default silent logger by a sub-logger (component: server)
	ftpServer.Logger = &server.LoggerAdapter{logger.WithField("component", "server")}

	// Preparing the SIGTERM handling
	go signalHandler()

	// Blocking call, behaving similarly to the http.ListenAndServe
	if onlyConf {
		logger.Warn("Only creating conf")

		return
	}

	if err := ftpServer.ListenAndServe(); err != nil {
		logger.Error("Problem listening", "err", err)
	}

	// We wait at most 1 minutes for all clients to disconnect
	if err := driver.WaitGracefully(time.Minute); err != nil {
		ftpServer.Logger.Warn("Problem stopping server", "err", err)
	}
}

func stop() {
	driver.Stop()

	if err := ftpServer.Stop(); err != nil {
		ftpServer.Logger.Error("Problem stopping server", "err", err)
	}
}

func signalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)

	for {
		sig := <-ch

		if sig == syscall.SIGTERM {
			stop()

			break
		}
	}
}
