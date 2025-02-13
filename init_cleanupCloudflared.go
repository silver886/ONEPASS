package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func initCleanupCloudflared() {
	// Set logger prefix
	logInitCleanupCloudflared := log.WithField("prefix", "initCleanupCloudflared")

	// Compose config path
	logInitCleanupCloudflared.Infoln("Compose config path . . .")
	if configPath, err := homedir.Expand("~/.cloudflared"); err != nil {
		logInitCleanupCloudflared.WithError(err).Errorln("Cannot compose config path")
		exit(35)
	} else if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logInitCleanupCloudflared.WithError(err).Infoln("No need to cleanup")
		return
	} else if err != nil {
		logInitCleanupCloudflared.WithError(err).Errorln("Cannot find config path")
		exit(36)
	} else if err := filepath.Walk(configPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), "-token.lock") {
			logInitCleanupCloudflared.Debugln(path)
			if err := os.Remove(path); err != nil {
				logInitCleanupCloudflared.WithError(err).Errorln("Cannot remove previous lock file")
				exit(35)
			}
		}
		return nil
	}); err != nil {
		logInitCleanupCloudflared.WithError(err).Errorln("Cannot get configs")
		exit(35)
	}
}
