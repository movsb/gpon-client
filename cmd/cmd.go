package cmd

import (
	"fmt"
	"html/template"
	"os"
	"strconv"

	gponclient "github.com/movsb/gpon-client/client"
	"github.com/spf13/cobra"
)

var gwinfoTemplate = template.Must(template.New(`gwinfoTemplate`).Parse(`LAN IPv4: {{ .LANIPv4 }}
WAN IPv4: {{ .WANIPv4 }}
MAC     : {{ .MAC }}
`))

// Client ...
var Client *gponclient.GponClient

// RootCmd ...
var RootCmd *cobra.Command

func init() {
	rootCmd := &cobra.Command{
		Use:   os.Args[0],
		Short: `A GPON (Tiānyì Gateway) client used to modify router configurations`,
	}
	RootCmd = rootCmd
	gwinfoCmd := &cobra.Command{
		Use:   `gwinfo`,
		Short: `show gateway information`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			info := Client.GetGatewayInfo()
			gwinfoTemplate.Execute(os.Stdout, info)
		},
	}
	rootCmd.AddCommand(gwinfoCmd)
	portMapsCmd := &cobra.Command{
		Use:   `portmaps`,
		Short: `manage port mappings`,
	}
	rootCmd.AddCommand(portMapsCmd)
	portMapsListCmd := &cobra.Command{
		Use:   `list`,
		Short: `list all port mappings`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			rules := Client.ListPortMappings()
			fmtstr := "%-5v%-16v%-12v%-12v%-20v%-12v%-6v\n"
			fmt.Printf(fmtstr, "ID", "Name", "Protocol", "OuterPort", "InnerIP", "InnerPort", "Enable")
			fmt.Println("-----------------------------------------------------------------------------------")
			for i, r := range rules {
				fmt.Printf(fmtstr, i+1, r.Name, r.Protocol, r.OuterPort, r.InnerIP, r.InnerPort, r.Enable)
			}
		},
	}
	portMapsCmd.AddCommand(portMapsListCmd)
	portMapsCreateCmd := &cobra.Command{
		Use:     `create <name> <protocol> <outer-port> <inner-ip> <inner-port>`,
		Short:   `craete a port mapping`,
		Example: `create nginx TCP 8080 192.168.1.6 80`,
		Args:    cobra.ExactArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			protocol := args[1]
			outerPort, _ := strconv.Atoi(args[2])
			innerIP := args[3]
			innerPort, _ := strconv.Atoi(args[4])
			Client.CreatePortMapping(name, protocol, outerPort, innerIP, innerPort)
		},
	}
	portMapsCmd.AddCommand(portMapsCreateCmd)
	portMapsDeleteCmd := &cobra.Command{
		Use:   `delete <name>`,
		Short: `delete a port mapping`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Client.DeletePortMapping(args[0])
		},
	}
	portMapsCmd.AddCommand(portMapsDeleteCmd)
	portMapsEnableCmd := &cobra.Command{
		Use:     `enable <name>`,
		Short:   `enable a port mapping`,
		Example: `enable nginx`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Client.EnablePortMapping(args[0], true)
		},
	}
	portMapsCmd.AddCommand(portMapsEnableCmd)
	portMapsDisableCmd := &cobra.Command{
		Use:     `disable <name>`,
		Short:   `disable a port mapping`,
		Example: `disable nginx`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Client.EnablePortMapping(args[0], false)
		},
	}
	portMapsCmd.AddCommand(portMapsDisableCmd)
	devicesCmd := &cobra.Command{
		Use:   `devices`,
		Short: `manage devices`,
	}
	rootCmd.AddCommand(devicesCmd)
	devicesListCmd := &cobra.Command{
		Use:   `list`,
		Short: `list devices`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			devices := Client.ListDevices()
			if len(devices) == 0 {
				return
			}
			fmtstr := "%-10v%-7v%-14v%15v%18v    %-10v%-10v%-14v\n"
			fmt.Printf(fmtstr, `Name`, `Wired`, `IPv4`, `Upload Speed`, `Download Speed`, `Type`, `System`, `MAC`)
			fmt.Println("----------------------------------------------------------------------------------------------------")
			for _, d := range devices {
				fmt.Printf(fmtstr, d.Name, d.Wired, d.IP,
					friendlySpeed(d.UploadSpeed), friendlySpeed(d.DownloadSpeed),
					d.Type, d.System, d.MAC,
				)
			}
		},
	}
	devicesCmd.AddCommand(devicesListCmd)
}
