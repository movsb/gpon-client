package main

import "fmt"

func main() {
	client := MustDial(
		defEnv("IP", "192.168.1.1"),
		defEnv("USERNAME", "useradmin"),
		defEnv("PASSWORD", ""),
	)
	portMappings := client.ListPortMappings()
	for i, p := range portMappings {
		fmt.Printf("%d: %v\n", i+1, p)
	}
}
