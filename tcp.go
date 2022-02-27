package main

import (
	"errors"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"./socks"
)

//Create a socks5 proxy server
func socksLocal(addr, server string, shadow func(net.Conn) net.Conn) {
	tcpLocal(addr, server, shadow, func(c net.Conn) (socks.Addr, error) { return socks.Handshake(c) })
}

func tcpLocal(addr, server string, shadow func(net.Conn) net.Conn, getAddr func(net.Conn) (socks.Addr, error)) {
	l, err := net.Listen("tcp", addr) //监听本地IP
	if err != nil {
		//logf("failed to listen on %s: %v", addr, err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			//logf("failed to accept: %s", err)
			continue
		}

		go func() {
			defer c.Close()
			tgt, err := getAddr(c) //获取目的端IP
			if err != nil {
				// logf("failed to get target address: %v", err)
				return
			}

			rc, err := net.Dial("tcp", server) //向服务端发送数据
			if err != nil {
				// logf("failed to connect to server %v: %v", server, err)
				return
			}
			defer rc.Close()
			// if config.TCPCork {
			// 	rc = timedCork(rc, 10*time.Millisecond, 1280)
			// }
			// rc = shadow(rc)

			if _, err = rc.Write(tgt); err != nil { //在数据流中写入目的端IP
				// logf("failed to send target address: %v", err)
				return
			}

			// logf("proxy %s <-> %s <-> %s", c.RemoteAddr(), server, tgt)
			if err = relay(rc, c); err != nil { //数据流Copy
				// logf("relay error: %v", err)
			}
		}()
	}
}

// Listen on addr for incoming connections.
// func tcpRemote(addr string, shadow func(net.Conn) net.Conn) {
// 	l, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		logf("failed to listen on %s: %v", addr, err)
// 		return
// 	}

// 	logf("listening TCP on %s", addr)
// 	for {
// 		c, err := l.Accept()
// 		if err != nil {
// 			logf("failed to accept: %v", err)
// 			continue
// 		}

// 		go func() {
// 			defer c.Close()
// 			if config.TCPCork {
// 				c = timedCork(c, 10*time.Millisecond, 1280)
// 			}
// 			sc := shadow(c)

// 			tgt, err := socks.ReadAddr(sc)
// 			if err != nil {
// 				logf("failed to get target address from %v: %v", c.RemoteAddr(), err)
// 				// drain c to avoid leaking server behavioral features
// 				// see https://www.ndss-symposium.org/ndss-paper/detecting-probe-resistant-proxies/
// 				_, err = io.Copy(ioutil.Discard, c)
// 				if err != nil {
// 					logf("discard error: %v", err)
// 				}
// 				return
// 			}

// 			rc, err := net.Dial("tcp", tgt.String())
// 			if err != nil {
// 				logf("failed to connect to target: %v", err)
// 				return
// 			}
// 			defer rc.Close()

// 			logf("proxy %s <-> %s", c.RemoteAddr(), tgt)
// 			if err = relay(sc, rc); err != nil {
// 				logf("relay error: %v", err)
// 			}
// 		}()
// 	}
// }

// relay copies between left and right bidirectionally
func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()
	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}
