package main

import (
	"strings"
	"sync"

	"github.com/silver886/execute"
	"github.com/sirupsen/logrus"
)

const (
	toolName    = "ONEPASS"
	toolVersion = "1.0.0"

	toolTitle = toolName + " " + toolVersion
)

func main() {
	// Set logger prefix
	logMain := log.WithField("prefix", "main")

	// Start init post
	logMain.Warnln("Keep this window to stay online")
	logMain.Warnln("Wait a moment while we set up the tunnel")
	logMain.Infoln("Start init post . . .")
	initPost()

	// Login
	logMain.Warnln("Log in to the tunnel on the page that just popped up")
	logMain.Infoln("Get the first remote . . .")
	firstRemote := arg.Remote
	if strings.Contains(firstRemote, separator) {
		firstRemote = strings.Split(firstRemote, separator)[0]
	}
	if cmd, err := execute.New(
		cloudflaredPath,
		"access",
		"login",
		firstRemote,
	).Hide().Run(); err != nil {
		logMain.WithError(err).WithFields(logrus.Fields{
			"cmd_path":   cmd.Path,
			"cmd_args":   cmd.Args,
			"cmd_stdout": cmd.OutString(),
			"cmd_stderr": cmd.ErrString(),
		}).Errorln("Cannot login")
	} else {
		logMain.WithFields(logrus.Fields{
			"hostname": arg.Remote,
			"stdout":   cmd.OutString(),
			"stderr":   cmd.ErrString(),
		}).Infoln("Logged in")
	}

	// Create tunnel
	logMain.Warnln("Setting up the tunnel")
	remotes, locals := strings.Split(arg.Remote, separator), strings.Split(arg.Local, separator)
	logMain.Infoln("Generate wait group . . .")
	var wg sync.WaitGroup
	wg.Add(len(remotes))
	for i := range remotes {
		remote := remotes[i]
		local := locals[i]
		go func() {
			if cmd, err := execute.New(
				cloudflaredPath,
				"access",
				"tcp",
				"--hostname",
				remote,
				"--url",
				local,
			).Hide().Start(); err != nil {
				logMain.WithError(err).WithFields(logrus.Fields{
					"cmd_path":   cmd.Path,
					"cmd_args":   cmd.Args,
					"cmd_stdout": cmd.OutString(),
					"cmd_stderr": cmd.ErrString(),
				}).Errorln("Cannot create tunnel")
			} else {
				logMain.Debugln("Wait for cloudflared exit . . .")
				cmd.Wait()
				logMain.WithError(err).WithFields(logrus.Fields{
					"hostname": arg.Remote,
					"url":      arg.Local,
					"stdout":   cmd.OutString(),
					"stderr":   cmd.ErrString(),
				}).Infoln("Tunnel terminated")
			}
			wg.Done()
		}()
	}

	// Halt here
	logMain.Warnln("Let's start")
	logMain.Infoln("Wait for tunnel terminated . . .")
	wg.Wait()

	// Exit program
	exit(0)
}
