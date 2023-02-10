package detect_ip

import (
	"net"
	"os"
	"sync"

	"github.com/chyroc/detect-ip/internal"
)

var isDebug = os.Getenv("DEBUG") != ""

func LockIPV4() net.IP {
	return detectIPByServers(false)
}

func LockIPV6() net.IP {
	return detectIPByServers(true)
}

func detectIPByServers(isV6 bool) net.IP {
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
				ip := internal.DetectIP(internal.HttpProdClient, apiList[nowIndex], isV6, isDebug)
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
