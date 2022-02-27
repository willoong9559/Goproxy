package cipher

import (
	"net"
	"strings"
)

type Cipher interface {
	StreamConnCipher
}

type StreamConnCipher interface {
	StreamConn(net.Conn) net.Conn
}

const (
	aeadAes128Gcm        = "AEAD_AES_128_GCM"
	aeadAes256Gcm        = "AEAD_AES_256_GCM"
	aeadChacha20Poly1305 = "AEAD_CHACHA20_POLY1305"
)

// List of AEAD ciphers: key size in bytes and constructor
var aeadList = map[string]struct {
	KeySize int
	New     func([]byte) (shadowaead.Cipher, error)
}{
	aeadAes128Gcm:        {16, shadowaead.AESGCM},
	aeadAes256Gcm:        {32, shadowaead.AESGCM},
	aeadChacha20Poly1305: {32, shadowaead.Chacha20Poly1305},
}

// PickCipher returns a Cipher of the given name. Derive key from password if given key is empty.
func PickCipher(name string, password string) (Cipher, error) {
	name = strings.ToUpper(name)

	switch name {
	case "CHACHA20-IETF-POLY1305":
		name = aeadChacha20Poly1305
	case "AES-128-GCM":
		name = aeadAes128Gcm
	case "AES-256-GCM":
		name = aeadAes256Gcm
	}

	if choice, ok := aeadList[name]; ok {
		return &aeadCipher{aead}, err
	}

	return nil, ErrCipherNotSupported
}
