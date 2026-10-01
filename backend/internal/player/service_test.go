package player

import "testing"

func TestCreateAndGet(t *testing.T) {
	s := NewService()
	p := Profile{ID: "p1", AccountID: "a1", DisplayName: "Test", Level: 1, Energy: 100}
	if err := s.Create(p); err != nil {
		t.Fatal(err)
	}
	got, ok := s.Get("p1")
	if !ok || got.DisplayName != "Test" {
		t.Fatal("profile was not stored")
	}
}
