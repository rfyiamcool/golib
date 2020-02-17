package dstparser

import (
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
)

const SO_ORIGINAL_DST = 80

var (
	ErrGetSocketoptIPv6 = errors.New("get socketopt ipv6 error")
	ErrResolveTCPAddr   = errors.New("resolve tcp address error")
	ErrTCPConn          = errors.New("not a valid TCPConn")
)

// For transparent proxy.
// Get REDIRECT package's originial dst address.
// Note: it may be only support linux.
func GetOriginalDstAddr(conn *net.TCPConn) (addr net.Addr, c *net.TCPConn, err error) {
	fc, errRet := conn.File()
	if errRet != nil {
		conn.Close()
		err = ErrTCPConn
		return
	} else {
		conn.Close()
	}
	defer fc.Close()

	mreq, errRet := syscall.GetsockoptIPv6Mreq(int(fc.Fd()), syscall.IPPROTO_IP, SO_ORIGINAL_DST)
	if errRet != nil {
		err = ErrGetSocketoptIPv6
		c, _ = getTCPConnFromFile(fc)
		return
	}

	// only support ipv4
	ip := net.IPv4(mreq.Multiaddr[4], mreq.Multiaddr[5], mreq.Multiaddr[6], mreq.Multiaddr[7])
	port := uint16(mreq.Multiaddr[2])<<8 + uint16(mreq.Multiaddr[3])
	addr, err = net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", ip.String(), port))
	if err != nil {
		err = ErrResolveTCPAddr
		return
	}

	c, errRet = getTCPConnFromFile(fc)
	if errRet != nil {
		err = ErrTCPConn
		return
	}
	return
}

func getTCPConnFromFile(f *os.File) (*net.TCPConn, error) {
	newConn, err := net.FileConn(f)
	if err != nil {
		return nil, ErrTCPConn
	}

	c, ok := newConn.(*net.TCPConn)
	if !ok {
		return nil, ErrTCPConn
	}
	return c, nil
}
