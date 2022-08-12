package hole

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"shiliyuanma/p2p-demo/x"
)

var (
	p2ps = make(map[string]*p2p)

	bufPoolUdp = sync.Pool{
		New: func() interface{} {
			return make([]byte, 2048)
		},
	}
)

type p2p struct {
	visitorAddr  *net.UDPAddr
	providerAddr *net.UDPAddr
}

func startRelay(port int) error {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{net.ParseIP("0.0.0.0"), port, ""})
	if err != nil {
		return err
	} else {
		fmt.Println(listener.LocalAddr().String())
	}

	defer func() { x.Hold = true }()
	go func() {
		for {
			select {
			case <-x.ExitC:
				return
			default:
			}

			buf := bufPoolUdp.Get().([]byte)
			n, addr, err := listener.ReadFromUDP(buf)
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				continue
			}

			go handleP2P(listener, addr, string(buf[:n]))
		}
	}()

	return nil
}

func handleP2P(listener *net.UDPConn, addr *net.UDPAddr, str string) {
	var (
		v  *p2p
		ok bool
	)
	arr := strings.Split(str, CONN_DATA_SEQ)
	arr = stringSliceDistinct(arr...)
	if len(arr) < 2 {
		return
	}
	if v, ok = p2ps[arr[0]]; !ok {
		v = new(p2p)
		p2ps[arr[0]] = v
	}
	fmt.Printf("new p2p connection ,role %s , password %s ,local address %s\n", arr[1], arr[0], addr.String())
	if arr[1] == WORK_P2P_VISITOR {
		v.visitorAddr = addr
		for i := 20; i > 0; i-- {
			if v.providerAddr != nil {
				listener.WriteTo([]byte(v.providerAddr.String()), v.visitorAddr)
				listener.WriteTo([]byte(v.visitorAddr.String()), v.providerAddr)
				break
			}
			time.Sleep(time.Second)
		}
		delete(p2ps, arr[0])
	} else {
		v.providerAddr = addr
	}
}
