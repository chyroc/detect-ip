package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
)

var (
	regxIPv4 = regexp.MustCompile(`(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`)
	regxIPv6 = regexp.MustCompile(`([0-9A-Fa-f]{0,4}:){2,7}([0-9A-Fa-f]{1,4}$|((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
)

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
		ipString = regxIPv6.FindString(body)
	} else {
		ipString = regxIPv4.FindString(body)
	}
	if isDebug {
		fmt.Printf("'%s' return ip=%s\n", apiServer, ipString)
	}
	return net.ParseIP(ipString)
}

func fetchURL(httpCli *http.Client, url string) (int, string, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "curl/7.79.1")
	res, err := httpCli.Do(req)
	if err != nil {
		return 0, "", err
	}
	bs, _ := io.ReadAll(res.Body)
	text := strings.TrimSpace(string(bs))

	if strings.HasPrefix(text, "{") {
		ipJSON := new(ipJSON)
		_ = json.Unmarshal(bs, ipJSON)
		if ipJSON.IP != "" {
			text = ipJSON.IP
		}
	}

	return res.StatusCode, text, nil
}

type ipJSON struct {
	IP string `json:"ip"`
}
