package inventory

import "testing"

func TestInventoryStackingAndRemoval(t *testing.T){
 s:=NewService();if err:=s.RegisterItem(Item{ID:"water",Name:"Water",Category:"supply",Stackable:true,MaxStack:10});err!=nil{t.Fatal(err)}
 if err:=s.Add("c1","water",15);err!=nil{t.Fatal(err)}
 inv:=s.Get("c1");if len(inv.Stacks)!=2||inv.Stacks[0].Quantity!=10||inv.Stacks[1].Quantity!=5{t.Fatalf("unexpected inventory: %+v",inv)}
 if err:=s.Remove("c1","water",12);err!=nil{t.Fatal(err)}
 inv=s.Get("c1");if len(inv.Stacks)!=1||inv.Stacks[0].Quantity!=3{t.Fatalf("unexpected after removal: %+v",inv)}
}
