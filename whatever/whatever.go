package whatever

import (
	"crypto/cipher"

	"./aead"
)

func NewWhateverClient(listen, server, psk string) (*SnellClient, error) {
	cipher := aead.NewAES128GCM([]byte(psk))
	sc := &SnellClient{
		server: server,
		cipher: cipher,
	}

	p, err := newSnellPool(MaxPoolCap, PoolTimeoutMS, sc.newSession)
	if err != nil {
		return nil, err
	}
	sc.pool = p

	sl, err := socks5.NewSocksProxy(listen, sc.handleSnell)
	if err != nil {
		return nil, err
	}
	sc.socks5 = sl

	return sc, nil
}
