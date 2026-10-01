package store

import (
	"context"
	"github.com/shido-ui/Worldbuster/backend/internal/simulation"
)

type EconomyBalanceRepository struct{ DB *DB }

type BalanceCycleResult struct {
	ProcessedAssets int
	PriceChanges    int
	TotalSupply     int64
	TotalDemand     int64
}

func (r EconomyBalanceRepository) Rebalance(ctx context.Context, limit int) (BalanceCycleResult, error) {
	if limit <= 0 {
		limit = 250
	}
	tx, err := r.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return BalanceCycleResult{}, err
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.QueryContext(ctx, `SELECT id::text,base_price,current_price,supply,demand FROM market_assets ORDER BY updated_at LIMIT $1 FOR UPDATE`, limit)
	if err != nil {
		return BalanceCycleResult{}, err
	}
	defer func() { _ = rows.Close() }()
	out := BalanceCycleResult{}
	for rows.Next() {
		var id string
		var s simulation.EconomySignal
		if err := rows.Scan(&id, &s.BasePrice, &s.CurrentPrice, &s.Supply, &s.Demand); err != nil {
			return out, err
		}
		out.ProcessedAssets++
		out.TotalSupply += s.Supply
		out.TotalDemand += s.Demand
		a := simulation.CalculatePriceAdjustment(s)
		if a.NewPrice == s.CurrentPrice {
			continue
		}
		if _, err := tx.ExecContext(ctx, `UPDATE market_assets SET current_price=$1,updated_at=NOW() WHERE id=$2`, a.NewPrice, id); err != nil {
			return out, err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO economy_price_history(asset_id,old_price,new_price,supply,demand,reason) VALUES($1,$2,$3,$4,$5,$6)`, id, s.CurrentPrice, a.NewPrice, s.Supply, s.Demand, "SUPPLY_DEMAND_REBALANCE"); err != nil {
			return out, err
		}
		out.PriceChanges++
		if _, err := tx.ExecContext(ctx, `UPDATE market_assets SET demand=GREATEST(0,(demand*9)/10),supply=GREATEST(0,(supply*9)/10) WHERE id=$1`, id); err != nil {
			return out, err
		}
	}
	if err := rows.Err(); err != nil {
		return out, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO economy_balance_cycles(processed_assets,price_changes,total_supply,total_demand) VALUES($1,$2,$3,$4)`, out.ProcessedAssets, out.PriceChanges, out.TotalSupply, out.TotalDemand); err != nil {
		return out, err
	}
	if err := tx.Commit(); err != nil {
		return out, err
	}
	return out, nil
}
