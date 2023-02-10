package detect_ip

import (
	_ "embed"
	"strings"
)

//go:embed ipv4.txt
var ipv4 string

//go:embed ipv6.txt
var ipv6 string

var ipv4ApiList []string

var ipv6ApiList []string

func init() {
	for _, v := range strings.Split(ipv4, "\n") {
		v = strings.TrimSpace(v)
		if v != "" {
			ipv4ApiList = append(ipv4ApiList, v)
		}
	}

	for _, v := range strings.Split(ipv6, "\n") {
		v = strings.TrimSpace(v)
		if v != "" {
			ipv6ApiList = append(ipv6ApiList, v)
		}
	}
}
