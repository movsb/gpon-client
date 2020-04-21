package client

import "net"

// RetVal ...
type RetVal struct {
	RetVal int `json:"retVal"`
}

// PortMappingRule ...
type PortMappingRule struct {
	Name      string `json:"desp"`
	Protocol  string `json:"protocol"`
	OuterPort int    `json:"exPort"`
	InnerIP   string `json:"client"`
	InnerPort int    `json:"inPort"`
	Enable    int    `json:"enable"`
}

// GatewayInfo ...
type GatewayInfo struct {
	LANIPv4 net.IP
	MAC     string
	WANIPv4 net.IP
}

type _GatewayInfo struct {
	LANIPv4 string `json:"LANIP"`
	MAC     string
	WANIPv4 string `json:"WANIP"`
}

// ToGatewayInfo ...
func (gw _GatewayInfo) ToGatewayInfo() GatewayInfo {
	info := GatewayInfo{
		LANIPv4: net.ParseIP(gw.LANIPv4),
		MAC:     gw.MAC,
		WANIPv4: net.ParseIP(gw.WANIPv4),
	}
	return info
}
