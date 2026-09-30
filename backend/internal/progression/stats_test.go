package progression

import "testing"

func TestStatBlock(t *testing.T){
 s:=NewStatBlock()
 if s.Total()!=5 { t.Fatal("unexpected base total") }
 if !s.Add("strength",4) || s.Strength!=5 { t.Fatal("stat increase failed") }
 if s.Add("unknown",1) { t.Fatal("unknown stat accepted") }
}
