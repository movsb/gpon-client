package main

type PortMappingRule struct {
	Enable      int    `json:"enable"`
	Protocol    string `json:"protocol"`
	InPort      int    `json:"inPort"`
	ExPort      int    `json:"exPort"`
	Description string `json:"desc"`
	ClientIP    string `json:"client"`
}
