package main

import (
	"os"

	"github.com/movsb/gpon-client/client"
	"github.com/movsb/gpon-client/cmd"
)

func defEnv(name string, def string) string {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	return value
}

func main() {
	cmd.Client = client.MustDial(
		defEnv("IP", "192.168.1.1"),
		defEnv("USERNAME", "useradmin"),
		defEnv("PASSWORD", ""),
	)
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
