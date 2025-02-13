package main

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	cloudflaredPath string
)

func initDownloadCloudflared() {
	// Set logger prefix
	logDownloadCloudflared := log.WithField("prefix", "downloadCloudflared")

	// Compose URL
	logDownloadCloudflared.Infoln("Compose URL . . .")
	url := bytes.NewBufferString("https://github.com/cloudflare/cloudflared/releases/")
	if arg.CloudflaredVersion == "latest" {
		if _, err := url.WriteString(arg.CloudflaredVersion); err != nil {
			logDownloadCloudflared.WithError(err).Errorln("Cannot compose URL")
			exit(35)
		} else if _, err := url.WriteString("/download"); err != nil {
			logDownloadCloudflared.WithError(err).Errorln("Cannot compose URL")
			exit(35)
		}
	} else {
		if _, err := url.WriteString("download/"); err != nil {
			logDownloadCloudflared.WithError(err).Errorln("Cannot compose URL")
			exit(35)
		} else if _, err := url.WriteString(arg.CloudflaredVersion); err != nil {
			logDownloadCloudflared.WithError(err).Errorln("Cannot compose URL")
			exit(35)
		}
	}

	// Download file
	logDownloadCloudflared.Infoln("Download file . . .")
	if _, err := url.WriteString("/cloudflared-windows-amd64.exe"); err != nil {
		logDownloadCloudflared.WithError(err).Errorln("Cannot compose URL")
		exit(35)
	} else if response, err := http.Get(url.String()); err != nil {
		logDownloadCloudflared.WithError(err).Errorln("Cannot send GET request")
		exit(35)
	} else if response.StatusCode != http.StatusOK {
		logDownloadCloudflared.WithFields(logrus.Fields{
			"url":        url.String(),
			"statusCode": response.StatusCode,
		}).Errorln("Cannot download file")
		exit(35)
	} else if tempFile, err := os.CreateTemp(workDir, "cloudflared*.exe"); err != nil {
		logDownloadCloudflared.WithError(err).Errorln("Cannot create file")
		exit(35)
	} else if _, err = io.Copy(tempFile, response.Body); err != nil {
		logDownloadCloudflared.WithError(err).Errorln("Cannot write file")
		exit(35)
	} else if err = response.Body.Close(); err != nil {
		logDownloadCloudflared.WithError(err).Errorln("Cannot close request")
		exit(35)
	} else if err = tempFile.Close(); err != nil {
		logDownloadCloudflared.WithError(err).Errorln("Cannot close file")
		exit(35)
	} else {
		cloudflaredPath = tempFile.Name()
	}
}
