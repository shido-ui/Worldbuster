package progression

import "testing"

func TestPlayerProgressionStore(t *testing.T){
 s:=NewPlayerProgressionStore()
 p,g:=s.AddXP("p1",100)
 if g!=1 || p.Level!=2 {t.Fatal("xp was not persisted")}
 if got:=s.Get("p1");got.Level!=2 {t.Fatal("progression was not retained")}
}
