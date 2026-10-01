package cache

import (
	"testing"
)

func TestNewRequiresAddress(t *testing.T) {
	if _, err := New("", "", 0); err == nil {
		t.Fatal("expected missing address error")
	}
}

func TestNewCreatesClient(t *testing.T) {
	client, err := New("localhost:6379", "", 0)
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("expected client")
	}
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
}
