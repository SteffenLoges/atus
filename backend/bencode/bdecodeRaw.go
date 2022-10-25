package bencode

import (
	"errors"

	"github.com/zeebo/bencode"
)

type DictRaw map[string]interface{}

func BDecodeRaw(data []byte) (DictRaw, error) {

	var dict DictRaw
	err := bencode.DecodeBytes(data, &dict)
	if err != nil {
		return nil, err
	}

	if _, ok := dict["info"]; !ok {
		return nil, errors.New("missing 'info' in dict")
	}

	return dict, nil

}
