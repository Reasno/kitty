package ots3

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestUploader(t *testing.T) {
	url := "https://via.placeholder.com/150"
	orig, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer orig.Body.Close()

	uploader := NewManager("Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", "https://play.minio.io:9000", "asia", "asiatrip")
	nu, err := uploader.Upload(context.Background(), orig.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(nu)
}
