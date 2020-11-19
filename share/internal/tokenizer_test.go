package internal

import (
	"fmt"
	"testing"
)

func TestTokenizer_Encode(t *testing.T) {
	h := NewTokenizer("foo")
	d, err := h.Encode(1452140)
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 10 {
		t.Fatal("length of data should be 10")
	}
	fmt.Println(d)
}

func TestTokenizer_Decode(t *testing.T) {
	h := NewTokenizer("foo")
	d := "bnJaEZ8WJN"
	i, err := h.Decode(d)
	if err != nil {
		t.Fatal(err)
	}
	if i != 1452140 {
		t.Fatal("i should be 1452140")
	}
	fmt.Println(i)
}
