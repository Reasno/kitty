package invitecode

import (
	"errors"

	"github.com/speps/go-hashids"
)

var ErrFailedToDecodeToken = errors.New("cannot decode token")

type Tokenizer struct {
	h *hashids.HashID
}

type TokenizerFactory func(salt string) *Tokenizer

func NewTokenizer(salt string) *Tokenizer {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 10
	h, _ := hashids.NewWithData(hd)
	return &Tokenizer{
		h,
	}
}

func (t Tokenizer) Decode(encoded string) (uint, error) {
	ints, err := t.h.DecodeWithError(encoded)
	if len(ints) > 0 {
		return uint(ints[0]), err
	}
	return 0, ErrFailedToDecodeToken
}

func (t Tokenizer) Encode(id uint) (string, error) {
	return t.h.Encode([]int{int(id)})
}
