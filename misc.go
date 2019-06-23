package main

import "os"

func defEnv(name string, def string) string {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	return value
}
