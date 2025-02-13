package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	cloudflaredPath string
)

func initDownloadCloudflared() {
	// Set logger prefix
	logInitDownloadCloudflared := log.WithField("prefix", "initDownloadCloudflared")

	// Compose URL
	logInitDownloadCloudflared.Infoln("Compose URL . . .")
	url := bytes.NewBufferString("https://github.com/cloudflare/cloudflared/releases/")
	if arg.CloudflaredVersion == "latest" {
		if _, err := url.WriteString(arg.CloudflaredVersion); err != nil {
			logInitDownloadCloudflared.WithError(err).Errorln("Cannot compose URL")
			exit(35)
		} else if _, err := url.WriteString("/download"); err != nil {
			logInitDownloadCloudflared.WithError(err).Errorln("Cannot compose URL")
			exit(35)
		}
	} else {
		if _, err := url.WriteString("download/"); err != nil {
			logInitDownloadCloudflared.WithError(err).Errorln("Cannot compose URL")
			exit(35)
		} else if _, err := url.WriteString(arg.CloudflaredVersion); err != nil {
			logInitDownloadCloudflared.WithError(err).Errorln("Cannot compose URL")
			exit(35)
		}
	}

	// Download file
	logInitDownloadCloudflared.Infoln("Download file . . .")
	if _, err := url.WriteString("/" + cloudflaredName); err != nil {
		logInitDownloadCloudflared.WithError(err).Errorln("Cannot compose URL")
		exit(35)
	} else if response, err := http.Get(url.String()); err != nil {
		logInitDownloadCloudflared.WithError(err).Errorln("Cannot send GET request")
		exit(35)
	} else if response.StatusCode != http.StatusOK {
		logInitDownloadCloudflared.WithFields(logrus.Fields{
			"url":        url.String(),
			"statusCode": response.StatusCode,
		}).Errorln("Cannot download file")
		exit(35)
	} else if tempFile, err := os.CreateTemp(workDir, "*."+strings.Split(cloudflaredName, ".")[strings.Count(cloudflaredName, ".")]); err != nil {
		logInitDownloadCloudflared.WithError(err).Errorln("Cannot create file")
		exit(35)
	} else if _, err = io.Copy(tempFile, response.Body); err != nil {
		logInitDownloadCloudflared.WithError(err).Errorln("Cannot write file")
		exit(35)
	} else if err = response.Body.Close(); err != nil {
		logInitDownloadCloudflared.WithError(err).Errorln("Cannot close request")
		exit(35)
	} else if err = tempFile.Close(); err != nil {
		logInitDownloadCloudflared.WithError(err).Errorln("Cannot close file")
		exit(35)
	} else {
		cloudflaredPath = tempFile.Name()
	}
}
