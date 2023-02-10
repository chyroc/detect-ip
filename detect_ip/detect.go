package detect_ip

import (
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

func LockIPV4() net.IP {
	return detectIP(ipv4ApiList, regxIPv4)
}

func LockIPV6() net.IP {
	return detectIP(ipv6ApiList, regxIPv6)
}

var regxIPv4 = regexp.MustCompile(`(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`)
var regxIPv6 = regexp.MustCompile(`([0-9A-Fa-f]{0,4}:){2,7}([0-9A-Fa-f]{1,4}$|((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)

func detectIP(apiList []string, reg *regexp.Regexp) net.IP {
	closed := make(chan struct{})
	var result net.IP = nil
	index := int32(-1)
	var safeClose = func() {
		select {
		case <-closed:
		default:
			close(closed)
		}
	}
	lock := sync.Mutex{}
	for i := 0; i < 5; i++ {
		go func() {
			select {
			case <-closed:
				return
			default:
				lock.Lock()
				index++
				nowIndex := index
				lock.Unlock()
				if nowIndex >= int32(len(apiList)) {
					return
				}
				statusCode, body, err := fetchURL(apiList[nowIndex])
				if err == nil && statusCode < 300 && statusCode >= 200 && body != "" {
					ipString := reg.FindString(strings.TrimSpace(body))
					if ip := net.ParseIP(ipString); ip != nil {
						lock.Lock()
						safeClose()
						result = ip
						lock.Unlock()
					}
				}
			}
		}()
	}
	return result
}

func fetchURL(url string) (int, string, error) {
	res, err := httpClient.Get(url)
	if err != nil {
		return 0, "", err
	}
	bs, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(bs), nil
}

var httpClient = &http.Client{
	Timeout: time.Second / 5,
}
