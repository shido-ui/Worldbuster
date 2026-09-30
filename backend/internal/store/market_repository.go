package store

import("context";"database/sql";"errors";"time")

var ErrMarketInvalid=errors.New("invalid market request")

type MarketAsset struct{ID string `json:"id"`;Symbol string `json:"symbol"`;Name string `json:"name"`;Category string `json:"category"`;BasePrice int64 `json:"basePrice"`;CurrentPrice int64 `json:"currentPrice"`;Supply int64 `json:"supply"`;Demand int64 `json:"demand"`;UpdatedAt time.Time `json:"updatedAt"`}
type MarketOrder struct{ID string `json:"id"`;AssetID string `json:"assetId"`;SellerAccountID string `json:"sellerAccountId"`;Quantity int64 `json:"quantity"`;UnitPrice int64 `json:"unitPrice"`;Status string `json:"status"`;CreatedAt time.Time `json:"createdAt"`}
type MarketRepository struct{DB *DB}

func(r MarketRepository) Assets(ctx context.Context)([]MarketAsset,error){rows,err:=r.DB.SQL.QueryContext(ctx,"SELECT id::text,symbol,name,category,base_price,current_price,supply,demand,updated_at FROM market_assets ORDER BY symbol");if err!=nil{return nil,err};defer rows.Close();out:=[]MarketAsset{};for rows.Next(){var x MarketAsset;if err:=rows.Scan(&x.ID,&x.Symbol,&x.Name,&x.Category,&x.BasePrice,&x.CurrentPrice,&x.Supply,&x.Demand,&x.UpdatedAt);err!=nil{return nil,err};out=append(out,x)};return out,rows.Err()}
func(r MarketRepository) Orders(ctx context.Context,assetID string)([]MarketOrder,error){if assetID==""{return nil,ErrMarketInvalid};rows,err:=r.DB.SQL.QueryContext(ctx,"SELECT id::text,asset_id::text,seller_account_id::text,quantity,unit_price,status,created_at FROM market_orders WHERE asset_id=$1 AND status='OPEN' ORDER BY unit_price,created_at LIMIT 100",assetID);if err!=nil{return nil,err};defer rows.Close();out:=[]MarketOrder{};for rows.Next(){var x MarketOrder;if err:=rows.Scan(&x.ID,&x.AssetID,&x.SellerAccountID,&x.Quantity,&x.UnitPrice,&x.Status,&x.CreatedAt);err!=nil{return nil,err};out=append(out,x)};return out,rows.Err()}

func(r MarketRepository) PlaceSell(ctx context.Context,accountID,assetID string,quantity,unitPrice int64)(MarketOrder,error){
 if accountID==""||assetID==""||quantity<=0||unitPrice<=0{return MarketOrder{},ErrMarketInvalid}
 tx,err:=r.DB.SQL.BeginTx(ctx,nil);if err!=nil{return MarketOrder{},err};defer tx.Rollback()
 var itemID sql.NullString
 if err=tx.QueryRowContext(ctx,"SELECT item_id FROM market_assets WHERE id=$1",assetID).Scan(&itemID);err!=nil{return MarketOrder{},err}
 if !itemID.Valid||itemID.String==""{return MarketOrder{},ErrMarketInvalid}
 var playerID string
 if err=tx.QueryRowContext(ctx,"SELECT id::text FROM player_profiles WHERE account_id=$1 FOR UPDATE",accountID).Scan(&playerID);err!=nil{return MarketOrder{},err}
 var qty int64
 if err=tx.QueryRowContext(ctx,"SELECT quantity FROM inventory_stacks WHERE player_id=$1 AND item_id=$2 FOR UPDATE",playerID,itemID.String).Scan(&qty);err!=nil||qty<quantity{return MarketOrder{},errors.New("insufficient inventory")}
 if _,err=tx.ExecContext(ctx,"UPDATE inventory_stacks SET quantity=quantity-$1 WHERE player_id=$2 AND item_id=$3",quantity,playerID,itemID.String);err!=nil{return MarketOrder{},err}
 if _,err=tx.ExecContext(ctx,"DELETE FROM inventory_stacks WHERE player_id=$1 AND item_id=$2 AND quantity<=0",playerID,itemID.String);err!=nil{return MarketOrder{},err}
 var x MarketOrder
 err=tx.QueryRowContext(ctx,"INSERT INTO market_orders(asset_id,seller_account_id,quantity,unit_price) VALUES($1,$2,$3,$4) RETURNING id::text,asset_id::text,seller_account_id::text,quantity,unit_price,status,created_at",assetID,accountID,quantity,unitPrice).Scan(&x.ID,&x.AssetID,&x.SellerAccountID,&x.Quantity,&x.UnitPrice,&x.Status,&x.CreatedAt)
 if err!=nil{return MarketOrder{},err}
 if err=tx.Commit();err!=nil{return MarketOrder{},err}
 return x,nil
}

func(r MarketRepository) Buy(ctx context.Context,buyer,orderID string,quantity int64)(MarketOrder,error){
 if buyer==""||orderID==""||quantity<=0{return MarketOrder{},ErrMarketInvalid}
 tx,err:=r.DB.SQL.BeginTx(ctx,nil);if err!=nil{return MarketOrder{},err};defer tx.Rollback()
 var x MarketOrder
 err=tx.QueryRowContext(ctx,"SELECT id::text,asset_id::text,seller_account_id::text,quantity,unit_price,status,created_at FROM market_orders WHERE id=$1 FOR UPDATE",orderID).Scan(&x.ID,&x.AssetID,&x.SellerAccountID,&x.Quantity,&x.UnitPrice,&x.Status,&x.CreatedAt)
 if err!=nil{return MarketOrder{},err}
 if x.Status!="OPEN"||x.SellerAccountID==buyer||quantity>x.Quantity{return MarketOrder{},ErrMarketInvalid}
 var itemID string
 if err=tx.QueryRowContext(ctx,"SELECT item_id FROM market_assets WHERE id=$1 AND item_id IS NOT NULL AND item_id<>''",x.AssetID).Scan(&itemID);err!=nil{return MarketOrder{},ErrMarketInvalid}
 if x.UnitPrice<=0||quantity>9223372036854775807/x.UnitPrice{return MarketOrder{},ErrMarketInvalid}
 total:=quantity*x.UnitPrice
 var balance int64;var buyerAccount string
 if err=tx.QueryRowContext(ctx,"SELECT id::text,balance FROM economy_accounts WHERE owner_account_id=$1 FOR UPDATE",buyer).Scan(&buyerAccount,&balance);err!=nil||balance<total{return MarketOrder{},errors.New("insufficient funds")}
 var sellerAccount string
 if err=tx.QueryRowContext(ctx,"SELECT id::text FROM economy_accounts WHERE owner_account_id=$1 FOR UPDATE",x.SellerAccountID).Scan(&sellerAccount);err!=nil{return MarketOrder{},err}
 if _,err=tx.ExecContext(ctx,"UPDATE economy_accounts SET balance=balance-$1,updated_at=NOW() WHERE id=$2",total,buyerAccount);err!=nil{return MarketOrder{},err}
 var buyerBalance int64; if err=tx.QueryRowContext(ctx,"SELECT balance FROM economy_accounts WHERE id=$1",buyerAccount).Scan(&buyerBalance);err!=nil{return MarketOrder{},err}
 if _,err=tx.ExecContext(ctx,"UPDATE economy_accounts SET balance=balance+$1,updated_at=NOW() WHERE id=$2",total,sellerAccount);err!=nil{return MarketOrder{},err}
 var sellerBalance int64; if err=tx.QueryRowContext(ctx,"SELECT balance FROM economy_accounts WHERE id=$1",sellerAccount).Scan(&sellerBalance);err!=nil{return MarketOrder{},err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO economy_ledger(id,account_id,amount,balance_after,reason,reference_id) VALUES(gen_random_uuid(),$1,$2,$3,'market_purchase',$4)",buyerAccount,-total,buyerBalance,x.ID);err!=nil{return MarketOrder{},err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO economy_ledger(id,account_id,amount,balance_after,reason,reference_id) VALUES(gen_random_uuid(),$1,$2,$3,'market_sale',$4)",sellerAccount,total,sellerBalance,x.ID);err!=nil{return MarketOrder{},err}
 var buyerPlayer string
 if err=tx.QueryRowContext(ctx,"SELECT id::text FROM player_profiles WHERE account_id=$1 FOR UPDATE",buyer).Scan(&buyerPlayer);err!=nil{return MarketOrder{},err}
var held int64
rows,err:=tx.QueryContext(ctx,"SELECT quantity FROM inventory_stacks WHERE player_id=$1 FOR UPDATE",buyerPlayer)
if err!=nil{return MarketOrder{},err}
for rows.Next(){var q int64;if err:=rows.Scan(&q);err!=nil{rows.Close();return MarketOrder{},err};held+=q}
if err:=rows.Close();err!=nil{return MarketOrder{},err}
if held>100||quantity>100-held{return MarketOrder{},errors.New("inventory capacity exceeded")}
 if _,err=tx.ExecContext(ctx,"INSERT INTO inventory_stacks(player_id,item_id,quantity) VALUES($1,$2,$3) ON CONFLICT(player_id,item_id) DO UPDATE SET quantity=inventory_stacks.quantity+EXCLUDED.quantity",buyerPlayer,itemID,quantity);err!=nil{return MarketOrder{},err}
 remaining:=x.Quantity-quantity;status:="OPEN";if remaining==0{status="FILLED"}
 if _,err=tx.ExecContext(ctx,"UPDATE market_orders SET quantity=$1,status=$2 WHERE id=$3",remaining,status,x.ID);err!=nil{return MarketOrder{},err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO market_trades(asset_id,buyer_account_id,seller_account_id,quantity,unit_price) VALUES($1,$2,$3,$4,$5)",x.AssetID,buyer,x.SellerAccountID,quantity,x.UnitPrice);err!=nil{return MarketOrder{},err}
 if _,err=tx.ExecContext(ctx,"UPDATE market_assets SET demand=demand+$2,current_price=GREATEST(1,((current_price*9)+$3)/10),updated_at=NOW() WHERE id=$1",x.AssetID,quantity,x.UnitPrice);err!=nil{return MarketOrder{},err}
 if err=tx.Commit();err!=nil{return MarketOrder{},err}
 x.Quantity=remaining;x.Status=status;return x,nil
}

func(r MarketRepository) Cancel(ctx context.Context,accountID,orderID string)(MarketOrder,error){
 if accountID==""||orderID==""{return MarketOrder{},ErrMarketInvalid}
 tx,err:=r.DB.SQL.BeginTx(ctx,nil);if err!=nil{return MarketOrder{},err};defer tx.Rollback()
 var x MarketOrder
 if err=tx.QueryRowContext(ctx,"SELECT id::text,asset_id::text,seller_account_id::text,quantity,unit_price,status,created_at FROM market_orders WHERE id=$1 FOR UPDATE",orderID).Scan(&x.ID,&x.AssetID,&x.SellerAccountID,&x.Quantity,&x.UnitPrice,&x.Status,&x.CreatedAt);err!=nil{return MarketOrder{},err}
 if x.Status!="OPEN"||x.SellerAccountID!=accountID{return MarketOrder{},ErrMarketInvalid}
 var itemID string; if err=tx.QueryRowContext(ctx,"SELECT item_id FROM market_assets WHERE id=$1 AND item_id IS NOT NULL AND item_id<>''",x.AssetID).Scan(&itemID);err!=nil{return MarketOrder{},ErrMarketInvalid}
 var playerID string; if err=tx.QueryRowContext(ctx,"SELECT id::text FROM player_profiles WHERE account_id=$1 FOR UPDATE",accountID).Scan(&playerID);err!=nil{return MarketOrder{},err}
 if _,err=tx.ExecContext(ctx,"INSERT INTO inventory_stacks(player_id,item_id,quantity) VALUES($1,$2,$3) ON CONFLICT(player_id,item_id) DO UPDATE SET quantity=inventory_stacks.quantity+EXCLUDED.quantity",playerID,itemID,x.Quantity);err!=nil{return MarketOrder{},err}
 if _,err=tx.ExecContext(ctx,"UPDATE market_orders SET status='CANCELLED',quantity=0 WHERE id=$1",orderID);err!=nil{return MarketOrder{},err}
 if err=tx.Commit();err!=nil{return MarketOrder{},err};x.Status="CANCELLED";x.Quantity=0;return x,nil
}
