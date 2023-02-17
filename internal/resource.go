package internal

import (
	_ "embed"
	"sort"
	"strconv"
	"strings"
)

// https://www.sitelike.org/similar/ipify.org/

//go:embed ipv4.txt
var ipv4 string

//go:embed ipv6.txt
var ipv6 string

type ApiServer struct {
	URL  string
	Ping int64
}

type ApiServerList []*ApiServer

func (a ApiServerList) Len() int {
	return len(a)
}

func (a ApiServerList) Less(i, j int) bool {
	if a[i].Ping == 0 {
		if a[j].Ping == 0 {
			return false // j < i
		}
		return false // j < i
	}
	if a[j].Ping == 0 {
		return true // i < j
	}
	return a[i].Ping < a[j].Ping
}

func (a ApiServerList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ApiServerList) URLs() []string {
	var r []string
	for _, v := range a {
		r = append(r, v.URL)
	}
	return r
}

func (a ApiServerList) ToText() string {
	s := new(strings.Builder)
	for _, v := range a {
		s.WriteString(v.URL)
		s.WriteString(",")
		s.WriteString(strconv.FormatInt(v.Ping, 10))
		s.WriteString("\n")
	}
	return s.String()
}

var (
	IpV4ApiList ApiServerList
	IpV6ApiList ApiServerList
)

func parseApiServer(s string) *ApiServer {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	l := strings.SplitN(s, ",", 2)
	if len(l) != 2 {
		return nil
	}
	if strings.TrimSpace(l[1]) == "" {
		return &ApiServer{URL: l[0], Ping: 0}
	}
	ping, err := strconv.ParseInt(strings.TrimSpace(l[1]), 10, 64)
	if err != nil {
		return &ApiServer{URL: l[0], Ping: ping}
	}
	return &ApiServer{URL: l[0], Ping: ping}
}

func init() {
	for _, v := range strings.Split(ipv4, "\n") {
		if api := parseApiServer(v); api != nil {
			IpV4ApiList = append(IpV4ApiList, api)
		}
	}

	for _, v := range strings.Split(ipv6, "\n") {
		if api := parseApiServer(v); api != nil {
			IpV6ApiList = append(IpV6ApiList, api)
		}
	}

	sort.Sort(IpV4ApiList)
	sort.Sort(IpV6ApiList)
}
