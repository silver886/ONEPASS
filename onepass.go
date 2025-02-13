package main

import (
	"strings"
	"sync"

	"github.com/silver886/execute"
	"github.com/silver886/logger"
	"github.com/sirupsen/logrus"
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
	logMain.Warnln("Log in to the tunnel")
	logMain.Infoln("Get the first remote . . .")
	firstRemote := arg.Remote
	if strings.Contains(firstRemote, separator) {
		firstRemote = strings.Split(firstRemote, separator)[0]
	}
	var loginWaitGroup sync.WaitGroup
	loginWaitGroup.Add(1)
	if cmd, err := execute.New(
		cloudflaredPath,
		"access",
		"login",
		"--quiet",
		firstRemote,
	).Hide().Start(); err != nil {
		logMain.WithError(err).WithFields(logrus.Fields{
			"cmd_path":   cmd.Path,
			"cmd_args":   cmd.Args,
			"cmd_stdout": cmd.OutString(),
			"cmd_stderr": cmd.ErrString(),
		}).Errorln("Cannot login")
		exit(1234)
	} else {
		go func(log *logger.Entry, cmd *execute.Cmd, remote string) {
			// Set logger prefix
			logLogin := log.WithField("sub", "login")

			// Get login URL
			logLogin.Infoln("Get login URL . . .")
			for msg := cmd.ErrStringNext(); cmd.ProcessState == nil; msg = cmd.ErrStringNext() {
				if msg != "" {
					for _, v := range strings.Split(msg, "\n") {
						if strings.Contains(v, "https://"+remote) {
							logLogin.WithFields(logrus.Fields{
								"url": v,
							}).Warnln("If the browser failed to open, visit the URL below")
						}
					}
				}
			}
		}(logMain, cmd, firstRemote)
		go func() {
			cmd.Wait()
			loginWaitGroup.Done()
		}()
		loginWaitGroup.Wait()
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
	var tunnelWaitGroup sync.WaitGroup
	tunnelWaitGroup.Add(len(remotes))
	for i := range remotes {
		go func(log *logger.Entry, remote string, local string) {
			// Set logger prefix
			logTunnel := log.WithFields(logrus.Fields{
				"sub":    "tunnel",
				"remote": remote,
				"local":  local,
			})

			// Create tunnel
			logTunnel.Infoln("Create tunnel . . .")
			if cmd, err := execute.New(
				cloudflaredPath,
				"access",
				"tcp",
				"--hostname",
				remote,
				"--url",
				local,
			).Hide().Start(); err != nil {
				logTunnel.WithError(err).WithFields(logrus.Fields{
					"cmd_path":   cmd.Path,
					"cmd_args":   cmd.Args,
					"cmd_stdout": cmd.OutString(),
					"cmd_stderr": cmd.ErrString(),
				}).Errorln("Cannot create tunnel")
			} else {
				logTunnel.Debugln("Wait for cloudflared exit . . .")
				cmd.Wait()
				logTunnel.WithError(err).WithFields(logrus.Fields{
					"hostname": arg.Remote,
					"url":      arg.Local,
					"stdout":   cmd.OutString(),
					"stderr":   cmd.ErrString(),
				}).Infoln("Tunnel terminated")
			}
			tunnelWaitGroup.Done()
		}(logMain, remotes[i], locals[i])
	}

	// Halt here
	logMain.Warnln("Let's start")
	logMain.Infoln("Wait for tunnel terminated . . .")
	tunnelWaitGroup.Wait()

	// Exit program
	exit(0)
}
