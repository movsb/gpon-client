package main

type RetVal struct {
	RetVal int `json:"retVal"`
}

type PortMappingRule struct {
	Name      string `json:"desp"`
	Protocol  string `json:"protocol"`
	OuterPort int    `json:"exPort"`
	InnerIP   string `json:"client"`
	InnerPort int    `json:"inPort"`
	Enable    int    `json:"enable"`
}
