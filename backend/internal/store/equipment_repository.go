package store

import ("context";"database/sql";"errors";"time")

type ItemDefinition struct {ID,Name,Category,Slot string; Stackable bool; MaxStack,PowerBonus,DefenseBonus,SpeedBonus,DurabilityMax,RequiredLevel int}
type EquipmentRecord struct {PlayerID,Slot,ItemID string; Durability int; EquippedAt time.Time}
type ItemProgressionRecord struct {PlayerID,ItemID string; UpgradeLevel int; Experience int64}
var ErrItemNotFound=errors.New("item not found")
var ErrEquipmentRequirement=errors.New("equipment requirement not met")

type EquipmentRepository struct{DB *DB}

func(r EquipmentRepository) ListItems(ctx context.Context)([]ItemDefinition,error){
 rows,err:=r.DB.SQL.QueryContext(ctx,"SELECT id,name,category,slot,stackable,max_stack,power_bonus,defense_bonus,speed_bonus,durability_max,required_level FROM item_definitions WHERE active ORDER BY category,name");if err!=nil{return nil,err};defer rows.Close();out:=[]ItemDefinition{}
 for rows.Next(){var x ItemDefinition;if err:=rows.Scan(&x.ID,&x.Name,&x.Category,&x.Slot,&x.Stackable,&x.MaxStack,&x.PowerBonus,&x.DefenseBonus,&x.SpeedBonus,&x.DurabilityMax,&x.RequiredLevel);err!=nil{return nil,err};out=append(out,x)};return out,rows.Err()
}

func(r EquipmentRepository) Equip(ctx context.Context,playerID,itemID string,level int)(EquipmentRecord,error){
 tx,err:=r.DB.SQL.BeginTx(ctx,nil);if err!=nil{return EquipmentRecord{},err};defer tx.Rollback()
 var slot string;var maxDur,req int
 if err=tx.QueryRowContext(ctx,"SELECT slot,durability_max,required_level FROM item_definitions WHERE id=$1 AND active",itemID).Scan(&slot,&maxDur,&req);err==sql.ErrNoRows{return EquipmentRecord{},ErrItemNotFound};if err!=nil{return EquipmentRecord{},err};if level<req{return EquipmentRecord{},ErrEquipmentRequirement}
 var qty int
 if err=tx.QueryRowContext(ctx,"SELECT quantity FROM inventory_stacks WHERE player_id=$1 AND item_id=$2 FOR UPDATE",playerID,itemID).Scan(&qty);err!=nil||qty<1{return EquipmentRecord{},ErrItemNotFound}
 if _,err=tx.ExecContext(ctx,"UPDATE inventory_stacks SET quantity=quantity-1 WHERE player_id=$1 AND item_id=$2",playerID,itemID);err!=nil{return EquipmentRecord{},err}
 if _,err=tx.ExecContext(ctx,"DELETE FROM inventory_stacks WHERE player_id=$1 AND item_id=$2 AND quantity<=0",playerID,itemID);err!=nil{return EquipmentRecord{},err}
 var x EquipmentRecord
 err=tx.QueryRowContext(ctx,"INSERT INTO player_equipment(player_id,slot,item_id,durability) VALUES($1,$2,$3,$4) ON CONFLICT(player_id,slot) DO UPDATE SET item_id=EXCLUDED.item_id,durability=EXCLUDED.durability,equipped_at=NOW() RETURNING player_id::text,slot,item_id,durability,equipped_at",playerID,slot,itemID,maxDur).Scan(&x.PlayerID,&x.Slot,&x.ItemID,&x.Durability,&x.EquippedAt)
 if err!=nil{return EquipmentRecord{},err};if err=tx.Commit();err!=nil{return EquipmentRecord{},err};return x,nil
}

func(r EquipmentRepository) ListEquipped(ctx context.Context,playerID string)([]EquipmentRecord,error){
 rows,err:=r.DB.SQL.QueryContext(ctx,"SELECT player_id::text,slot,item_id,durability,equipped_at FROM player_equipment WHERE player_id=$1 ORDER BY slot",playerID);if err!=nil{return nil,err};defer rows.Close();out:=[]EquipmentRecord{}
 for rows.Next(){var x EquipmentRecord;if err:=rows.Scan(&x.PlayerID,&x.Slot,&x.ItemID,&x.Durability,&x.EquippedAt);err!=nil{return nil,err};out=append(out,x)};return out,rows.Err()
}
