package cmd

import "fmt"

var speedNames = []struct {
	Base int
	Name string
}{
	{1 << 30, "GB/s"},
	{1 << 20, "MB/s"},
	{1 << 10, "KB/s"},
	{1 << 0, "B/s"},
}

func friendlySpeed(b int) string {
	if b == 0 {
		return `0B/s`
	}
	var s string
	for i := 0; i < len(speedNames); i++ {
		n := speedNames[i]
		if b >= n.Base {
			s += fmt.Sprintf(`%.2f`, float64(b)/float64(n.Base))
			s += n.Name
			break
		}
	}
	return s
}
