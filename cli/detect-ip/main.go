package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/chyroc/detect-ip/detect_ip"
)

func main() {
	i := flag.Int64("timeout", 200, "timeout(ms)")
	flag.Parse()

	timeout := time.Duration(*i) * time.Millisecond

	start := time.Now()
	ipv4 := detect_ip.LockIPV4(detect_ip.WithTimeout(timeout))
	ipv4Ping := time.Since(start)
	ipv6 := detect_ip.LockIPV6(detect_ip.WithTimeout(timeout))
	ipv6Ping := time.Since(start) - ipv4Ping

	fmt.Printf("ipv4: %s, ping: %s\n", ipv4.String(), ipv4Ping.String())
	fmt.Printf("ipv6: %s, ping: %s\n", ipv6.String(), ipv6Ping.String())
}
