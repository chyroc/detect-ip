package detect_ip

import (
	"net"
	"sync"

	"github.com/chyroc/detect-ip/internal"
)

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
				ip := internal.DetectIP(internal.HttpProdClient, apiList[nowIndex], isV6)
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
	return result
}
