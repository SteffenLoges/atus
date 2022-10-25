package bencode

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	"github.com/zeebo/bencode"
)

func hashBuffer(buf bytes.Buffer) (string, error) {
	hash := sha1.New()
	hash.Write(buf.Bytes())
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (dr DictRaw) GenHash() (string, error) {
	var infoBuf bytes.Buffer
	bencode.NewEncoder(&infoBuf).Encode(dr["info"])
	return hashBuffer(infoBuf)
}

func (d *Dict) GenHash() (string, error) {
	var infoBuf bytes.Buffer
	bencode.NewEncoder(&infoBuf).Encode(d.Info)
	return hashBuffer(infoBuf)
}
