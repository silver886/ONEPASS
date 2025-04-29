package main

import (
	"os"
)

func initSetupCloudflared(path string) {
	// Set logger prefix
	logInitSetupCloudflared := log.WithField("prefix", "initSetupCloudflared")

	// Set file mode
	logInitSetupCloudflared.Infoln("Set file mode . . .")
	if err := os.Chmod(path, 0755); err != nil {
		logInitSetupCloudflared.WithError(err).Errorln("Cannot set file mode")
		exit(35)
	}
}
