package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"./whatever"
)

func main() {
	var flags struct {
		Client string
		Server string
		Socks  string
	}

	flag.StringVar(&flags.Server, "s", "", "server listen address or url")
	flag.StringVar(&flags.Client, "c", "", "client connect address or url")
	flag.StringVar(&flags.Socks, "socks", "", "listen address")
	flag.Parse()

	if flags.Server == "" && flags.Client == "" {
		flag.Usage()
		return
	}

	if flags.Client != "" {
		addr := flags.Client
		//ciphercode := "aes-128-gcm"
		var password string
		var err error

		//ciph, err := cipher.PickCipher(ciphercode, password)

		if strings.HasPrefix(addr, "whatever://") {
			addr, password, err = parseURL(addr)
			if err != nil {
				log.Fatal(err)
			}
		}
		sn, err := whatever.NewSnellClient(flags.Socks, addr, password)
		if err != nil {
			log.Fatalf("Failed to initialize snell client %v\n", err)
		}

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		sn.Close()
	}

	// Server
	// if flags.Server != "" {
	// 	addr := flags.Server
	// 	password := flags.Password
	// 	var err error

	// 	if strings.HasPrefix(addr, "whatever://") {
	// 		addr, password, err = parseURL(addr)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}

	// 	go tcpRemote(addr, password)
	// }
}

// parseURL parse url
func parseURL(s string) (addr, password string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}
	addr = u.Host
	password = u.User.Username()
	return
}
