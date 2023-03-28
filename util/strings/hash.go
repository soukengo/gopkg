package strings

import "github.com/speps/go-hashids/v2"

type HashID struct {
	hash *hashids.HashID
	err  error
}

func NewHashID(salt string, minLen int) *HashID {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLen
	hash, err := hashids.NewWithData(hd)
	return &HashID{hash: hash, err: err}
}

func (h *HashID) Encode(numbers []int64) (string, error) {
	return h.hash.EncodeInt64(numbers)
}
func (h *HashID) Decode(hash string) ([]int64, error) {
	return h.hash.DecodeInt64WithError(hash)
}
