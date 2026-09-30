package organization
import "testing"
func TestOrganizationMembership(t *testing.T){s:=NewService();if err:=s.Create(Organization{ID:"c1",Name:"Aster Works",Type:TypeCompany,OwnerID:"p1",MaxMembers:2});err!=nil{t.Fatal(err)};if err:=s.Join("c1","p1","owner");err!=nil{t.Fatal(err)};if err:=s.Join("c1","p2","member");err!=nil{t.Fatal(err)};if err:=s.Join("c1","p3","member");err!=ErrFull{t.Fatalf("got %v",err)};_,m,_:=s.Get("c1");if len(m)!=2{t.Fatal("membership missing")}}
