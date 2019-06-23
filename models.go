package main

type RetVal struct {
	RetVal int `json:"retVal"`
}

type PortMappingRule struct {
	Enable      int    `json:"enable"`
	Protocol    string `json:"protocol"`
	InnerPort   int    `json:"inPort"`
	OuterPort   int    `json:"exPort"`
	Description string `json:"desp"`
	InnerIP     string `json:"client"`
}
