package economy
import "testing"
func TestLedger(t *testing.T){s:=NewService();s.EnsureAccount("a");if _,err:=s.Credit("a",1000,"quest","q1");err!=nil{t.Fatal(err)};if _,err:=s.Debit("a",250,"purchase","p1");err!=nil{t.Fatal(err)};a,_:=s.Balance("a");if a.Balance!=750{t.Fatalf("balance=%d",a.Balance)};if _,err:=s.Debit("a",1000,"purchase","p2");err!=ErrInsufficientFunds{t.Fatalf("expected insufficient funds, got %v",err)};if len(s.History("a"))!=2{t.Fatal("ledger entries missing")}}
