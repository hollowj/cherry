package cherryNet

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

var localIPv4Str = "0.0.0.0"
var localIPv4Once = new(sync.Once)

func LocalIPV4() string {
	localIPv4Once.Do(func() {
		if ias, err := net.InterfaceAddrs(); err == nil {
			for _, address := range ias {
				if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						localIPv4Str = ipNet.IP.String()
						return
					}
				}
			}
		}
	})
	return localIPv4Str
}

func GetIPV4(addr net.Addr) string {
	if addr == nil {
		return ""
	}

	if ipNet, ok := addr.(*net.TCPAddr); ok {
		return ipNet.IP.String()
	}

	if ipNet, ok := addr.(*net.UDPAddr); ok {
		return ipNet.IP.String()
	}

	return ""
}

func LocalIPWithAutoPortForGame(min int32, max int32) string {
	ip := LocalIPV4()
	rand.Seed(time.Now().UnixNano())
	for {
		site := fmt.Sprintf("%s:%d", ip, min+rand.Int31n(max-min))
		listen, err := net.Listen("tcp", site)
		if err == nil {
			err = listen.Close()
			if err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			return site
		}
		time.Sleep(time.Millisecond * 100)
	}
}
