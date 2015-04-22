// proxy socks5 spec , pls refer http://blog.chinaunix.net/uid-26548237-id-3434356.html
// github project, pls refer to https://github.com/eahydra/socks
package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type ClientConn struct {
	conn net.Conn
}

func NewClientConn(conn net.Conn) (*ClientConn, error) {
	clientConn := &ClientConn{
		conn: conn,
	}
	return clientConn, nil
}

func (c *ClientConn) Read(data []byte) (int, error) {
	return c.conn.Read(data)
}

func (c *ClientConn) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}
func (c *ClientConn) Close() error {
	return c.conn.Close()
}
func (c *ClientConn) Run() {
	defer c.Close()

	// Socks step 1 : handshake
	if err := handShake(c); err != nil {
		fmt.Println("handshake failed, err:", err)
		return
	}

	// socks step 2 : CMD
	cmd, destHost, destPort, _, err := getCommand(c)
	if err != nil {
		fmt.Println("getCommand failed, err:", err)
		return
	}

	// response
	reply := []byte{0x05, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x22, 0x22}
	if cmd != 0x01 {
		reply[1] = 0x07 // unsupported command
		c.Write(reply)
		return
	}

	var dest io.ReadWriteCloser
	destConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", destHost, destPort))
	if err != nil {
		reply[1] = 0x05
		c.Write(reply)
		return
	}
	defer destConn.Close()
	dest = destConn

	reply[1] = 0x00
	if _, err = c.Write(reply); err != nil {
		return
	}
	fmt.Println("Socks5: step 2 <- ", reply)

	if dest == nil {
		panic("dest is nil")
	}

	fmt.Println("Socks5: step 3 < - > ")
	go func() {
		io.Copy(dest, c)
		fmt.Println("Client")
	}()
	io.Copy(c, dest)
	fmt.Println("Serve")

}

func handShake(rw io.ReadWriter) error {
	// version(1)+numMethods(1)+[256]methods
	buff := make([]byte, 258)
	n, err := io.ReadAtLeast(rw, buff, 2)
	if err != nil {
		return err
	}

	// version
	if buff[0] != 5 {
		return errors.New("socks unsuported version")
	}

	// numMethods
	numMethod := int(buff[1])

	// methods
	numMethod += 2
	if n <= numMethod {
		if _, err = io.ReadFull(rw, buff[n:numMethod]); err != nil {
			return err
		}
	}
	fmt.Println("Socks5: Step 1 -> ", buff[:numMethod])

	// return data
	buff[1] = 0 // no authentication
	if _, err := rw.Write(buff[:2]); err != nil {
		return err
	}
	fmt.Println("Socks5: Step 1 <- ", buff[:2])

	return nil
}

func getCommand(reader io.Reader) (cmd byte, destHost string, destPort uint16, data []byte, err error) {
	// version(1) + cmd(1) + reserved(1) + addrType(1) + domainLength(1) + maxDomainLength(256) + port(2)
	buff := make([]byte, 263)
	var n int
	n, err = io.ReadAtLeast(reader, buff, 9)
	if err != nil {
		return
	}

	// version
	if buff[0] != 5 {
		err = errors.New("socks unsuported version")
		return
	}

	// cmd
	cmd = buff[1]

	// addrType
	totalLength := 0
	switch buff[3] {
	case 0x01: // ipv4
		totalLength = 3 + 1 + 4 + 2 // veriosn + cmd + reserved + addrType + ip + 2
	case 0x03: // Domain
		totalLength = 3 + 1 + 1 + int(buff[4]) + 2 // veriosn + cmd + reserved + addrType + domainLength + lenght + 2
	case 0x04: // ipv6
		totalLength = 3 + 1 + 16 + 2 // veriosn + cmd + reserved + addrType + ipv6 + 2
	}
	if n < totalLength {
		if _, err = io.ReadFull(reader, buff[n:totalLength]); err != nil {
			return
		}
	}

	switch buff[3] {
	case 0x01:
		destHost = net.IP(buff[4 : 4+net.IPv4len]).String()
	case 0x03:
		destHost = string(buff[5 : 5+int(buff[4])])
	case 0x4:
		destHost = net.IP(buff[4 : 4+net.IPv6len]).String()
	}
	destPort = binary.BigEndian.Uint16(buff[totalLength-2 : totalLength])
	data = buff[:totalLength]
	fmt.Println("Socks5: Step 2 -> ", buff[:totalLength])
	return

}

func main() {
	// create server to listen
	listener, err := net.Listen("tcp", "127.0.0.1:8989")
	if err != nil {
		fmt.Println("can't create listener", err)
		return
	}
	defer listener.Close()

	for {
		// server get data
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("can't get data :", err)
			return
		}

		// connect remote server
		fmt.Println("Connect remote server")
		if clientConn, err := NewClientConn(conn); err == nil {
			go clientConn.Run()
		} else {
			fmt.Println("New Client fail: ", err)
		}

	}

}
