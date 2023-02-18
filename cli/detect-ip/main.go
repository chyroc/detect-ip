package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/chyroc/detect-ip/detect_ip"
)

func main() {
	i := flag.Int64("timeout", 5000, "timeout(ms)")
	isIPv4 := flag.Bool("v4", true, "print ipv4")
	isIPv6 := flag.Bool("v6", true, "print ipv6")
	isJson := flag.Bool("json", false, "print json")
	flag.Parse()

	result := detectIP(*i, *isIPv4, *isIPv6)

	result.Print(*isJson)
}

type Result struct {
	IPv4     string `json:"ipv4,omitempty"`
	IPv4Ping string `json:"ipv4_ping,omitempty"`
	IPv6     string `json:"ipv6,omitempty"`
	IPv6Ping string `json:"ipv6_ping,omitempty"`
}

func detectIP(i int64, isIPv4, isIPv6 bool) *Result {
	timeout := time.Duration(i) * time.Millisecond
	result := new(Result)

	start := time.Now()
	ipv4Ping := time.Duration(0)
	if isIPv4 {
		ipv4 := detect_ip.PublicIPV4(detect_ip.WithTimeout(timeout))
		result.IPv4 = ipv4.String()
		ipv4Ping = time.Since(start)
		result.IPv4Ping = ipv4Ping.String()
	}
	if isIPv6 {
		ipv6 := detect_ip.PublicIPV6(detect_ip.WithTimeout(timeout))
		result.IPv6 = ipv6.String()
		result.IPv6Ping = (time.Since(start) - ipv4Ping).String()
	}

	return result
}

func (r Result) Print(isJSON bool) {
	if isJSON {
		bs, _ := json.MarshalIndent(r, "", "  ")
		fmt.Println(string(bs))
	} else {
		if r.IPv4 != "" {
			fmt.Printf("ipv4: %s, ping: %s\n", r.IPv4, r.IPv4Ping)
		}
		if r.IPv6 != "" {
			fmt.Printf("ipv6: %s, ping: %s\n", r.IPv6, r.IPv6Ping)
		}
	}
}
