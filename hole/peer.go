package hole

import (
	"errors"
	"fmt"
	"net"
	"time"

	"shiliyuanma/p2p-demo/x"
)

func startPeer(relay, role, pwd string) error {
	fmt.Println("relay addr:", relay)
	relayAddr, err := net.ResolveUDPAddr("udp", relay)
	if err != nil {
		return err
	}
	localAddr, err := udpTempAddr()
	if err != nil {
		return err
	}
	fmt.Println("local addr:", localAddr)
	localConn, err := udpListen(localAddr)
	if err != nil {
		return err
	}

	var (
		buf        = make([]byte, 1024)
		remoteAddr string
	)
	for {
		localConn.SetWriteDeadline(time.Now().Add(time.Second * 5))
		if _, err := localConn.WriteTo(getWriteStr(pwd, role), relayAddr); err != nil {
			fmt.Println("write to server error:", err)
			continue
		}

		localConn.SetReadDeadline(time.Now().Add(time.Second * 5))
		if count, _, err := localConn.ReadFromUDP(buf); err != nil {
			fmt.Println("read from server error:", err)
			continue
		} else {
			remoteAddr = string(buf[:count])
		}
		if len(remoteAddr) > 0 {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Println("remote addr:", remoteAddr)
	rAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		return err
	}
	if err = holePunch(localConn, rAddr); err != nil {
		return err
	} else {
		fmt.Println("hole punch success")
	}

	defer func() { x.Hold = true }()
	if conn, err := udpListen(localAddr); err != nil {
		return err
	} else {
		go handlePeer(conn, rAddr, role)
	}

	return nil
}

func holePunch(localConn *net.UDPConn, rAddr *net.UDPAddr) error {
	defer localConn.Close()
	isClose := false
	defer func() { isClose = true }()

	go func() {
		fmt.Printf("try send test packet to target %s\n", rAddr)
		ticker := time.NewTicker(time.Millisecond * 500)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if isClose {
					return
				}
				localConn.SetWriteDeadline(time.Now().Add(time.Second * 5))
				if _, err := localConn.WriteTo([]byte(WORK_P2P_CONNECT), rAddr); err != nil {
					fmt.Println("write connect error:", err)
				}
			}
		}
	}()

	buf := make([]byte, 10)
	for {
		localConn.SetReadDeadline(time.Now().Add(time.Second * 5))
		n, addr, err := localConn.ReadFromUDP(buf)
		localConn.SetReadDeadline(time.Time{})
		if err != nil {
			break
		}
		switch string(buf[:n]) {
		case WORK_P2P_SUCCESS:
			for i := 20; i > 0; i-- {
				if _, err = localConn.WriteTo([]byte(WORK_P2P_END), addr); err != nil {
					return err
				}
			}
			return nil
		case WORK_P2P_END:
			fmt.Printf("Remotely Address %s Reply Packet Successfully Received\n", addr.String())
			return nil
		case WORK_P2P_CONNECT:
			go func() {
				for i := 20; i > 0; i-- {
					fmt.Printf("try send receive success packet to target %s\n", addr.String())
					if _, err = localConn.WriteTo([]byte(WORK_P2P_SUCCESS), addr); err != nil {
						return
					}
					time.Sleep(time.Second)
				}
			}()
		}
	}

	return errors.New("connect to the target failed, maybe the nat type is not support p2p")
}

func handlePeer(conn *net.UDPConn, rAddr *net.UDPAddr, role string) {
	defer conn.Close()
	fmt.Println("start handle")
	for {
		select {
		case <-x.ExitC:
			return
		default:
		}

		if role == WORK_P2P_VISITOR {
			conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
			if _, err := conn.WriteTo([]byte("hello"), rAddr); err != nil {
				fmt.Println("write error:", err)
			} else {
				fmt.Println("write: hello")
			}
		} else if role == WORK_P2P_PROVIDER {
			data := make([]byte, 1024)
			conn.SetReadDeadline(time.Now().Add(time.Second * 5))
			if count, _, err := conn.ReadFromUDP(data); err != nil {
				fmt.Println("read error:", err)
			} else {
				fmt.Println("read:", string(data[:count]))
			}
		}
		time.Sleep(time.Second)
	}
}
