package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"
)

func main() {
	if len(os.Args) == 1 {
		usage(os.Args[0])
		return
	}
	client := MustDial(
		defEnv("IP", "192.168.1.1"),
		defEnv("USERNAME", "useradmin"),
		defEnv("PASSWORD", ""),
	)
	switch cmd := os.Args[1]; cmd {
	default:
		log.Fatalf("unknown command: %s\n", cmd)
	case "portmaps":
		subcmd := os.Args[2]
		a := os.Args[3:]
		_ = a
		switch subcmd {
		default:
			log.Fatalf("unknown sub-command: %s\n", subcmd)
		case "list":
			rules := client.ListPortMappings()
			fmtstr := "%-5v%-16v%-12v%-12v%-20v%-12v%-6v\n"
			fmt.Printf(fmtstr, "ID", "Name", "Protocol", "OuterPort", "InnerIP", "InnerPort", "Enable")
			fmt.Println("-----------------------------------------------------------------------------------")
			for i, r := range rules {
				fmt.Printf(fmtstr, i+1, r.Name, r.Protocol, r.OuterPort, r.InnerIP, r.InnerPort, r.Enable)
			}
		case "create":
			name := a[0]
			protocol := a[1]
			outerPort, _ := strconv.Atoi(a[2])
			innerIP := a[3]
			innerPort, _ := strconv.Atoi(a[4])
			client.CreatePortMapping(name, protocol, outerPort, innerIP, innerPort)
		case "delete":
			name := a[0]
			client.DeletePortMapping(name)
		case "enable":
			name := a[0]
			client.EnablePortMapping(name, true)
		case "disable":
			name := a[0]
			client.EnablePortMapping(name, false)
		}
	case `gwinfo`:
		info := client.GetGatewayInfo()
		gwinfoTemplate.Execute(os.Stdout, info)
	}
}

func usage(name string) {
	usageText := `A GPON (Tiānyì Gateway) client used to modify router configurations.

Usage: %[1]s command [sub-command] [<arguments>...]

All Available Command List:

	portmaps list
	portmaps create  <name> <protocol> <outer-port> <inner-ip> <inner-port>
	portmaps delete  <name>
	portmaps enable  <name>
	portmaps disable <name>

	gwinfo

`
	fmt.Fprintf(os.Stderr, usageText, name)
}

var gwinfoTemplate = template.Must(template.New(`gwinfoTemplate`).Parse(`LAN IPv4: {{ .LANIPv4 }}
WAN IPv4: {{ .WANIPv4 }}
MAC     : {{ .MAC }}
`))
