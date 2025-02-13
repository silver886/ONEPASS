package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func exitInit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	go func() {
		for sig := range c {
			log.WithFields(logrus.Fields{
				"prefix": "signal",
				"signal": sig,
			}).Infoln("Get signal")
			exit(7)
		}
	}()
}

func exit(code int) {
	// Set logger prefix
	logExit := log.WithField("prefix", "exit")

	logExit.WithField("exit_code", code).Infoln("Exiting program . . .")

	// Exit program with exit code
	logExit.WithField("exit_code", code).Infoln("Exit . . .")
	// Wait logger completed
	log.Wait()
	os.Exit(code)
}
