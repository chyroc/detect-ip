package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/chyroc/detect-ip/internal"
)

func main() {
	speedTestAndWrite(internal.IpV4ApiList, false, "./internal/ipv4.txt")
	speedTestAndWrite(internal.IpV6ApiList, true, "./internal/ipv6.txt")
}

func speedTestAndWrite(list internal.ApiServerList, isV6 bool, file string) {
	done := map[string]bool{}
	for _, v := range list {
		host, err := url.Parse(v.URL)
		if err != nil {
			log.Fatalln(err)
			return
		}
		if done[host.String()] {
			continue
		}
		done[host.String()] = true

		ip, ping := speedTest(v.URL, isV6)
		if ip == nil {
			v.Ping = 0
		} else {
			v.Ping = ping.Milliseconds()
		}
	}
	sort.Sort(list)

	err := ioutil.WriteFile(file, []byte(list.ToText()), 0o644)
	if err != nil {
		log.Fatalln(err)
	}
}

var httpTestClient = &http.Client{Timeout: time.Second * 5}

func speedTest(apiServer string, isV6 bool) (net.IP, time.Duration) {
	fmt.Printf("[test] start %s, v6=%v\n", apiServer, isV6)
	start := time.Now()
	ip := internal.DetectIP(httpTestClient, apiServer, isV6, false)
	dur := time.Since(start)
	fmt.Printf("[test] '%s' return %s, dur=%s\n", apiServer, ip.String(), dur.String())
	return ip, dur
}
