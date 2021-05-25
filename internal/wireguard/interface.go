package wireguard

import "io"

type Wireguard interface {
	io.Closer
	CreateConfigForNewKeys() (io.Reader, error)
	CreateConfigForPublicKey(key string) (io.Reader, error)
}
