package store

import (
	"context"
	"database/sql"
	"github.com/shido-ui/Worldbuster/backend/internal/simulation"
)

type TerritoryCycleResult struct {
	TerritoriesProcessed int
	ControlChanges       int
	StabilityChanges     int
}

func (r TerritoryRepository) ResolveConsequences(ctx context.Context, limit int) (TerritoryCycleResult, error) {
	if limit <= 0 {
		limit = 250
	}
	tx, err := r.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return TerritoryCycleResult{}, err
	}
	defer func() { _ = tx.Rollback() }()
	rows, err := tx.QueryContext(ctx, `SELECT id::text,controlling_organization_id::text,stability FROM territories ORDER BY updated_at LIMIT $1 FOR UPDATE`, limit)
	if err != nil {
		return TerritoryCycleResult{}, err
	}
	defer func() { _ = rows.Close() }()
	out := TerritoryCycleResult{}
	for rows.Next() {
		var id string
		var current *string
		var stability int
		if err := rows.Scan(&id, &current, &stability); err != nil {
			return out, err
		}
		var topOrg string
		var topVal, secondVal int
		err = tx.QueryRowContext(ctx, `SELECT organization_id::text,influence FROM territory_influence WHERE territory_id=$1 ORDER BY influence DESC,organization_id LIMIT 1`, id).Scan(&topOrg, &topVal)
		if err != nil && err != sql.ErrNoRows {
			return out, err
		}
		if topOrg != "" {
			if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(influence),0) FROM territory_influence WHERE territory_id=$1 AND organization_id<>$2`, id, topOrg).Scan(&secondVal); err != nil {
				return out, err
			}
		}
		outcome := simulation.ResolveTerritory(simulation.TerritoryPressure{TopInfluence: topVal, SecondInfluence: secondVal, Stability: stability})
		changed := false
		if topVal > 0 {
			if current == nil || *current != topOrg {
				_, err = tx.ExecContext(ctx, `UPDATE territories SET controlling_organization_id=$1,influence=$2,stability=$3,updated_at=NOW() WHERE id=$4`, topOrg, topVal, outcome.Stability, id)
				changed = true
			} else {
				_, err = tx.ExecContext(ctx, `UPDATE territories SET influence=$1,stability=$2,updated_at=NOW() WHERE id=$3`, topVal, outcome.Stability, id)
			}
		} else {
			_, err = tx.ExecContext(ctx, `UPDATE territories SET controlling_organization_id=NULL,influence=0,stability=$1,updated_at=NOW() WHERE id=$2`, outcome.Stability, id)
		}
		if err != nil {
			return out, err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO territory_effects(territory_id,commerce_modifier_bps,safety_modifier,stability_modifier,updated_at) VALUES($1,$2,$3,$4,NOW()) ON CONFLICT(territory_id) DO UPDATE SET commerce_modifier_bps=EXCLUDED.commerce_modifier_bps,safety_modifier=EXCLUDED.safety_modifier,stability_modifier=EXCLUDED.stability_modifier,updated_at=NOW()`, id, outcome.CommerceModifierBPS, outcome.SafetyModifier, outcome.Stability-stability); err != nil {
			return out, err
		}
		if changed {
			out.ControlChanges++
		}
		if outcome.Stability != stability {
			out.StabilityChanges++
		}
		out.TerritoriesProcessed++
	}
	if err := rows.Err(); err != nil {
		return out, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO territory_consequence_cycles(territories_processed,control_changes,stability_changes) VALUES($1,$2,$3)`, out.TerritoriesProcessed, out.ControlChanges, out.StabilityChanges); err != nil {
		return out, err
	}
	if err := tx.Commit(); err != nil {
		return out, err
	}
	return out, nil
}
