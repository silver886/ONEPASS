package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	workDir string
)

func initWorkdir() {
	// Set logger prefix
	logInitWorkdir := log.WithField("prefix", "initWorkdir")

	// Create workdir
	logInitWorkdir.Infoln("Create workdir . . .")
	if tempDir, err := os.MkdirTemp("", toolTitle); err != nil {
		logInitWorkdir.WithError(err).Errorln("Cannot create workdir")
		exit(35)
	} else {
		// Change workdir
		logInitWorkdir.Infoln("Change workdir . . .")
		logInitWorkdir.WithFields(logrus.Fields{
			"workdir": tempDir,
		}).Debugln("Change workdir . . .")
		if err := os.Chdir(tempDir); err != nil {
			logInitWorkdir.WithError(err).Errorln("Cannot change workdir")
			exit(36)
		}
		workDir = tempDir
	}
}
