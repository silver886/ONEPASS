package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

func initSetupCloudflared(archivePath string) {
	// Set logger prefix
	logInitSetupCloudflared := log.WithField("prefix", "initSetupCloudflared")

	// Open archive
	logInitSetupCloudflared.Infoln("Open archive . . .")
	if archive, err := os.Open(archivePath); err != nil {
		logInitSetupCloudflared.WithError(err).Errorln("Cannot open file")
		exit(667)
	} else if uncompressedStream, err := gzip.NewReader(archive); err != nil {
		logInitSetupCloudflared.WithError(err).Errorln("Cannot read archive")
		exit(667)
	} else {
		// Read content from archive
		logInitSetupCloudflared.Infoln("Read content from archive . . .")
		for tarReader := tar.NewReader(uncompressedStream); true; {
			if header, err := tarReader.Next(); err == io.EOF {
				break
			} else if err != nil {
				logInitSetupCloudflared.WithError(err).Errorln("Cannot read archive content")
				exit(667)
			} else if header.Typeflag == tar.TypeReg && header.Name == "cloudflared" {
				// Extract file from archive
				logInitSetupCloudflared.Infoln("Extract file from archive . . .")
				if tempFile, err := os.CreateTemp(workDir, ""); err != nil {
					logInitSetupCloudflared.WithError(err).Errorln("Cannot create file")
					exit(667)
				} else if _, err := io.Copy(tempFile, tarReader); err != nil {
					logInitSetupCloudflared.WithError(err).Errorln("Cannot write file")
					exit(667)
				} else if err = tempFile.Close(); err != nil {
					logInitSetupCloudflared.WithError(err).Errorln("Cannot close file")
					exit(35)
				} else if err := os.Chmod(tempFile.Name(), 0755); err != nil {
					logInitSetupCloudflared.WithError(err).Errorln("Cannot set file mode")
					exit(35)
				} else {
					cloudflaredPath = tempFile.Name()
				}
			}
		}
	}
}
