package main

import (
	"io/ioutil"
	"log"
	"sort"

	"github.com/chyroc/detect-ip/internal"
)

func main() {
	speedTestAndWrite(internal.IpV4ApiList, false, "./internal/ipv4.txt")
	speedTestAndWrite(internal.IpV6ApiList, true, "./internal/ipv6.txt")
}

func speedTestAndWrite(list internal.ApiServerList, isV6 bool, file string) {
	for _, v := range list {
		ip, ping := internal.SpeedTest(v.URL, isV6)
		if ip == nil {
			v.Ping = 0
		} else {
			v.Ping = ping.Milliseconds()
		}
	}
	sort.Sort(list)

	err := ioutil.WriteFile(file, []byte(list.ToText()), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
