package store

import (
 "context"
 "github.com/shido-ui/Worldbuster/backend/internal/inventory"
)

type InventoryRepository struct{DB *DB}

func(r InventoryRepository) Get(ctx context.Context,playerID string)(inventory.Inventory,error){
 rows,err:=r.DB.SQL.QueryContext(ctx,"SELECT item_id,quantity FROM inventory_stacks WHERE player_id=$1 ORDER BY item_id",playerID);if err!=nil{return inventory.Inventory{},err};defer rows.Close()
 out:=inventory.Inventory{CharacterID:playerID,Capacity:100}
 for rows.Next(){var s inventory.Stack;if err:=rows.Scan(&s.ItemID,&s.Quantity);err!=nil{return inventory.Inventory{},err};out.Stacks=append(out.Stacks,s)}
 return out,rows.Err()
}

func(r InventoryRepository) Add(ctx context.Context,playerID,itemID string,quantity int)error{
 if quantity<=0{return ErrInvalidQuantity}
 _,err:=r.DB.SQL.ExecContext(ctx,"INSERT INTO inventory_stacks(player_id,item_id,quantity) VALUES($1,$2,$3) ON CONFLICT(player_id,item_id) DO UPDATE SET quantity=inventory_stacks.quantity+EXCLUDED.quantity",playerID,itemID,quantity)
 return err
}
