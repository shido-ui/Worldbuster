package store

import (
 "context"
 "github.com/shido-ui/Worldbuster/backend/internal/simulation"
)

type ProductionRepository struct{DB *DB}

type ProductionCycleResult struct {
 Processed int
 Produced int
 InputsConsumed int64
 OutputsCreated int64
 LaborSpent int64
}

func(r ProductionRepository) SeedBusinessInventory(ctx context.Context,businessID,itemID string,quantity int64) error {
 if quantity<=0{return nil}
 _,err:=r.DB.SQL.ExecContext(ctx,`INSERT INTO simulated_business_inventory(business_id,item_id,quantity) VALUES($1,$2,$3)
 ON CONFLICT(business_id,item_id) DO UPDATE SET quantity=simulated_business_inventory.quantity+EXCLUDED.quantity,updated_at=NOW()`,businessID,itemID,quantity)
 return err
}

func(r ProductionRepository) RunCycle(ctx context.Context,limit int)(ProductionCycleResult,error){
 if limit<=0 {limit=100}
 tx,err:=r.DB.SQL.BeginTx(ctx,nil);if err!=nil{return ProductionCycleResult{},err};defer tx.Rollback()
 rows,err:=tx.QueryContext(ctx,`SELECT b.id::text,b.business_type,b.cash_balance,b.level,b.reputation,b.employees,b.revenue_total,b.expense_total,b.active,
  r.id::text,r.name,r.business_type,r.min_level,r.input_item,r.input_quantity,r.output_item,r.output_quantity,r.labor_cost
  FROM simulated_businesses b CROSS JOIN production_recipes r
  WHERE b.active=true AND r.active=true AND b.business_type=r.business_type AND b.level>=r.min_level
  ORDER BY b.updated_at,r.id LIMIT $1`,limit)
 if err!=nil{return ProductionCycleResult{},err};defer rows.Close()
 out:=ProductionCycleResult{}
 for rows.Next(){
  var b simulation.BusinessState;var r simulation.ProductionRecipe
  if err:=rows.Scan(&b.ID,&b.Type,&b.Cash,&b.Level,&b.Reputation,&b.Employees,&b.RevenueTotal,&b.ExpenseTotal,&b.Active,&r.ID,&r.Name,&r.BusinessType,&r.MinLevel,&r.InputItem,&r.InputQuantity,&r.OutputItem,&r.OutputQuantity,&r.LaborCost);err!=nil{return out,err}
  out.Processed++
  var input int64
  err=tx.QueryRowContext(ctx,`SELECT quantity FROM simulated_business_inventory WHERE business_id=$1 AND item_id=$2 FOR UPDATE`,b.ID,r.InputItem).Scan(&input)
  if err!=nil {continue}
  if !simulation.CanProduce(b,r,input){continue}
  _,err=tx.ExecContext(ctx,`UPDATE simulated_business_inventory SET quantity=quantity-$1,updated_at=NOW() WHERE business_id=$2 AND item_id=$3`,r.InputQuantity,b.ID,r.InputItem);if err!=nil{return out,err}
  _,err=tx.ExecContext(ctx,`INSERT INTO simulated_business_inventory(business_id,item_id,quantity) VALUES($1,$2,$3)
   ON CONFLICT(business_id,item_id) DO UPDATE SET quantity=simulated_business_inventory.quantity+EXCLUDED.quantity,updated_at=NOW()`,b.ID,r.OutputItem,r.OutputQuantity);if err!=nil{return out,err}
  _,err=tx.ExecContext(ctx,`UPDATE simulated_businesses SET cash_balance=cash_balance-$1,expense_total=expense_total+$1,updated_at=NOW() WHERE id=$2`,r.LaborCost,b.ID);if err!=nil{return out,err}
  _,err=tx.ExecContext(ctx,`INSERT INTO simulated_business_ledger(business_id,amount,balance_after,reason)
   SELECT $1,-$2,cash_balance,'PRODUCTION_LABOR' FROM simulated_businesses WHERE id=$1`,b.ID,r.LaborCost);if err!=nil{return out,err}
  _,err=tx.ExecContext(ctx,`INSERT INTO simulated_production_runs(business_id,recipe_id,input_item,input_quantity,output_item,output_quantity,labor_cost)
   VALUES($1,$2,$3,$4,$5,$6,$7)`,b.ID,r.ID,r.InputItem,r.InputQuantity,r.OutputItem,r.OutputQuantity,r.LaborCost);if err!=nil{return out,err}
  _,err=tx.ExecContext(ctx,`INSERT INTO world_supply_signals(item_id,supply,demand,updated_at) VALUES($1,$2,0,NOW())
   ON CONFLICT(item_id) DO UPDATE SET supply=world_supply_signals.supply+$2,updated_at=NOW()`,r.OutputItem,r.OutputQuantity);if err!=nil{return out,err}
  out.Produced++;out.InputsConsumed+=int64(r.InputQuantity);out.OutputsCreated+=int64(r.OutputQuantity);out.LaborSpent+=r.LaborCost
 }
 if err:=rows.Err();err!=nil{return out,err}
 if err:=tx.Commit();err!=nil{return out,err}
 return out,nil
}
