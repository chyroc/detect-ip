package internal

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

var regxIPv4 = regexp.MustCompile(`(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`)
var regxIPv6 = regexp.MustCompile(`([0-9A-Fa-f]{0,4}:){2,7}([0-9A-Fa-f]{1,4}$|((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)

func DetectIP(httpCli *http.Client, apiServer string, isV6 bool, isDebug bool) net.IP {
	statusCode, body, err := fetchURL(httpCli, apiServer)
	if isDebug {
		fmt.Printf("'%s' return %d, %s, err=%v\n", apiServer, statusCode, body, err)
	}
	if err != nil || statusCode < 200 || statusCode >= 300 || body == "" {
		return nil
	}
	ipString := ""
	if isV6 {
		ipString = regxIPv6.FindString(strings.TrimSpace(body))
	} else {
		ipString = regxIPv4.FindString(strings.TrimSpace(body))
	}
	if isDebug {
		fmt.Printf("'%s' return ip=%s\n", apiServer, ipString)
	}
	return net.ParseIP(ipString)
}

func fetchURL(httpCli *http.Client, url string) (int, string, error) {
	res, err := httpCli.Get(url)
	if err != nil {
		return 0, "", err
	}
	bs, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(bs), nil
}
