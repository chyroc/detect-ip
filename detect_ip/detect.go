package detect_ip

import (
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/chyroc/detect-ip/internal"
)

var isDebug = os.Getenv("DEBUG") != ""

type Option struct {
	timeout time.Duration
}

func WithTimeout(t time.Duration) func(o *Option) {
	return func(o *Option) {
		o.timeout = t
	}
}

func PublicIPV4(options ...func(o *Option)) net.IP {
	opt := makeOption(options)
	return detectIPByServers(false, opt.timeout)
}

func PublicIPV6(options ...func(o *Option)) net.IP {
	opt := makeOption(options)
	return detectIPByServers(true, opt.timeout)
}

func makeOption(options []func(o *Option)) *Option {
	o := &Option{
		timeout: time.Second / 5, // default timeout
	}
	for _, v := range options {
		v(o)
	}
	return o
}

func detectIPByServers(isV6 bool, timeout time.Duration) net.IP {
	apiList := []string{}
	if isV6 {
		apiList = internal.IpV6ApiList.URLs()
	} else {
		apiList = internal.IpV4ApiList.URLs()
	}

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
	wait := new(sync.WaitGroup)
	for i := 0; i < 5; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
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
				ip := internal.DetectIP(getHttpClient(timeout), apiList[nowIndex], isV6, isDebug)
				if ip == nil {
					return
				}

				lock.Lock()
				safeClose()
				result = ip
				lock.Unlock()
			}
		}()
	}
	wait.Wait()
	return result
}

var httpClientMap = sync.Map{}

func getHttpClient(timeout time.Duration) *http.Client {
	if v, ok := httpClientMap.Load(timeout); ok {
		return v.(*http.Client)
	}
	client := &http.Client{Timeout: timeout}
	httpClientMap.Store(timeout, client)
	return client
}
