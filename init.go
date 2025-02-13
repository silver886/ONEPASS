package main

import (
	"bytes"
	"flag"
	"strings"

	"github.com/rifflock/lfshook"
	"github.com/silver886/logger"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	separator string = ","
)

var (
	log      *logger.Logger
	debugLog *bytes.Buffer

	arg struct {
		Debug              bool
		CloudflaredVersion string
		Remote             string
		Local              string
	}

	defaultCloudflaredVersion string
	defaultRemote             string
	defaultLocal              string
)

func init() {
	// Define arguments
	flag.BoolVar(&arg.Debug, "debug", false, "Enable debug mode")
	flag.StringVar(&arg.CloudflaredVersion, "cloudflaredVersion", defaultCloudflaredVersion, "Cloudflared version")
	flag.StringVar(&arg.Remote, "remote", defaultRemote, "Remote host")
	flag.StringVar(&arg.Local, "local", defaultLocal, "Local listening port")

	// Create logger
	log, _ = logger.New(toolName,
		logrus.WarnLevel,
		[]logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
		false,
	)

	// Add debug log
	debugLog = bytes.NewBuffer([]byte{})
	log.Hooks.Add(lfshook.NewHook(
		debugLog,
		&prefixed.TextFormatter{FullTimestamp: true, ForceFormatting: true, DisableColors: true, SpacePadding: 64},
	))

	// Set logger prefix
	logInit := log.WithField("prefix", "init")

	logInit.Debugln("Create logger")
}

func initPost() {
	// Set logger prefix
	logInitPost := log.WithField("prefix", "initPost")

	// Parse arguments
	flag.Parse()

	if arg.Debug {
		log.Wait()
		log.Config(
			logrus.Level(len(logrus.AllLevels)-1),
			logrus.AllLevels,
			false,
		)
		logInitPost.Debugln("Update logger")
	}

	logInitPost.WithFields(logrus.Fields{
		"debug":               arg.Debug,
		"cloudflared_version": arg.CloudflaredVersion,
		"remote":              arg.Remote,
		"local":               arg.Local,
	}).Debugln("Command arguments")

	// Setup exit signal handler
	logInitPost.Infoln("Setup exit signal handler . . .")
	exitInit()

	// Check argument
	logInitPost.Infoln("Check argument . . .")
	if arg.CloudflaredVersion == "" {
		logInitPost.Errorln("Argument `cloudflaredVersion` should not be empty")
		exit(1)
	}
	if arg.Remote == "" {
		logInitPost.Errorln("Argument `remote` should not be empty")
		exit(1)
	}
	if arg.Local == "" {
		logInitPost.Errorln("Argument `local` should not be empty")
		exit(1)
	}
	if strings.Count(arg.Remote, separator) != strings.Count(arg.Local, separator) {
		logInitPost.Errorln("Argument `remote` and `local` should have same amount of elements")
		exit(1)
	}

	// Create working directory
	logInitPost.Infoln("Create working directory . . .")
	initWorkdir()

	// Download cloudflared
	logInitPost.Infoln("Download cloudflared . . .")
	initDownloadCloudflared()

	// Extract cloudflared
	logInitPost.Infoln("Extract cloudflared . . .")
	initExtractCloudflared(cloudflaredPath)

	// Cleanup cloudflared
	logInitPost.Infoln("Cleanup cloudflared . . .")
	initCleanupCloudflared()

	// Finish initialization
	logInitPost.Infoln("Finish initialization . . .")
}
