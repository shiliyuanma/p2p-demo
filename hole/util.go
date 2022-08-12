package hole

import (
	"bytes"
	"encoding/binary"
	"net"
)

func udpTempAddr() (string, error) {
	conn, err := net.Dial("udp", "114.114.114.114:53")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return conn.LocalAddr().String(), nil
}

func udpListen(addr string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	return udpConn, nil
}

// get seq str
func getWriteStr(v ...string) []byte {
	buffer := new(bytes.Buffer)
	var l int32
	for _, v := range v {
		l += int32(len([]byte(v))) + int32(len([]byte(CONN_DATA_SEQ)))
		binary.Write(buffer, binary.LittleEndian, []byte(v))
		binary.Write(buffer, binary.LittleEndian, []byte(CONN_DATA_SEQ))
	}
	return buffer.Bytes()
}

func stringSliceDistinct(arr ...string) []string {
	dict := make(map[string]string)
	out := make([]string, 0, 0)
	for _, i := range arr {
		dict[i] = i
	}
	for k, _ := range dict {
		out = append(out, k)
	}
	return out
}
