package otredis

import "fmt"

type KeyManager struct {
	Prefix string
}

func (k *KeyManager) Key(parts ...string) string {
	s := k.Prefix
	for _, part := range parts {
		s = fmt.Sprintf("%s:%s", s, part)
	}
	return s
}
func (k *KeyManager) Add(parts ...string) {
	for _, part := range parts {
		k.Prefix = fmt.Sprintf("%s:%s", k.Prefix, part)
	}
}
