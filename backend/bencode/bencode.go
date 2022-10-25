package bencode

import (
	"github.com/zeebo/bencode"
)

func (d *Dict) BEncode() ([]byte, error) {
	return bencode.EncodeBytes(d)
}
