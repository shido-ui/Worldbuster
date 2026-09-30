package store

import "testing"

func TestEquipmentRepositoryTypes(t *testing.T) {
 x:=ItemDefinition{ID:"field-boots",Slot:"FOOTWEAR",SpeedBonus:3}
 if x.ID==""||x.SpeedBonus!=3{t.Fatal("invalid item definition")}
}
