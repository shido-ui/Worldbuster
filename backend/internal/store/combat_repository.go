package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shido-ui/Worldbuster/backend/internal/simulation"
	"time"
)

var ErrCombatInvalid = errors.New("combat request invalid")
var ErrCombatCooldown = errors.New("combat cooldown active")
var ErrCombatEnergy = errors.New("insufficient energy")

type CombatRecord struct {
	ID, AttackerID, DefenderID               string
	WinnerID                                 *string
	AttackerScore, DefenderScore, EnergyCost int
	XPReward                                 int64
	CreatedAt                                time.Time
}

type CombatRepository struct{ DB *DB }

func (r CombatRepository) Resolve(ctx context.Context, attackerID, defenderID string) (CombatRecord, error) {
	if attackerID == "" || defenderID == "" || attackerID == defenderID {
		return CombatRecord{}, ErrCombatInvalid
	}
	tx, err := r.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return CombatRecord{}, err
	}
	defer func() { _ = tx.Rollback() }()
	var a, d struct {
		level, strength, defense, speed, endurance, energy int
		power, defenseBonus, speedBonus                    int
	}
	if err = tx.QueryRowContext(ctx, "SELECT level,strength,defense,speed,endurance,energy FROM player_profiles WHERE id=$1 FOR UPDATE", attackerID).Scan(&a.level, &a.strength, &a.defense, &a.speed, &a.endurance, &a.energy); err != nil {
		return CombatRecord{}, err
	}
	if err = tx.QueryRowContext(ctx, "SELECT COALESCE(SUM(i.power_bonus),0),COALESCE(SUM(i.defense_bonus),0),COALESCE(SUM(i.speed_bonus),0) FROM player_equipment e JOIN item_definitions i ON i.id=e.item_id WHERE e.player_id=$1", attackerID).Scan(&a.power, &a.defenseBonus, &a.speedBonus); err != nil {
		return CombatRecord{}, err
	}
	if err = tx.QueryRowContext(ctx, "SELECT level,strength,defense,speed,endurance,energy FROM player_profiles WHERE id=$1 FOR UPDATE", defenderID).Scan(&d.level, &d.strength, &d.defense, &d.speed, &d.endurance, &d.energy); err != nil {
		return CombatRecord{}, err
	}
	if err = tx.QueryRowContext(ctx, "SELECT COALESCE(SUM(i.power_bonus),0),COALESCE(SUM(i.defense_bonus),0),COALESCE(SUM(i.speed_bonus),0) FROM player_equipment e JOIN item_definitions i ON i.id=e.item_id WHERE e.player_id=$1", defenderID).Scan(&d.power, &d.defenseBonus, &d.speedBonus); err != nil {
		return CombatRecord{}, err
	}
	if a.energy < 10 {
		return CombatRecord{}, ErrCombatEnergy
	}
	var available time.Time
	err = tx.QueryRowContext(ctx, "SELECT available_at FROM combat_cooldowns WHERE player_id=$1", attackerID).Scan(&available)
	if err == nil && time.Now().UTC().Before(available) {
		return CombatRecord{}, ErrCombatCooldown
	}
	result := simulation.ResolveCombat(simulation.CombatInput{AttackerStrength: a.strength + a.power, AttackerSpeed: a.speed + a.speedBonus, AttackerDefense: a.defense + a.defenseBonus, AttackerEndurance: a.endurance, DefenderStrength: d.strength + d.power, DefenderSpeed: d.speed + d.speedBonus, DefenderDefense: d.defense + d.defenseBonus, DefenderEndurance: d.endurance, AttackerLevel: a.level, DefenderLevel: d.level})
	_, err = tx.ExecContext(ctx, "UPDATE player_profiles SET energy=energy-$2,updated_at=NOW() WHERE id=$1", attackerID, result.EnergyCost)
	if err != nil {
		return CombatRecord{}, err
	}
	availableAt := time.Now().UTC().Add(30 * time.Second)
	if _, err = tx.ExecContext(ctx, "INSERT INTO combat_cooldowns(player_id,available_at) VALUES($1,$2) ON CONFLICT(player_id) DO UPDATE SET available_at=EXCLUDED.available_at", attackerID, availableAt); err != nil {
		return CombatRecord{}, err
	}
	var winner interface{}
	if result.Winner == 1 {
		winner = attackerID
	} else if result.Winner == 2 {
		winner = defenderID
	}
	var id string
	if err = tx.QueryRowContext(ctx, "INSERT INTO combat_matches(attacker_id,defender_id,winner_id,attacker_score,defender_score,energy_cost,xp_reward) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id::text", attackerID, defenderID, winner, result.AttackerScore, result.DefenderScore, result.EnergyCost, result.XPReward).Scan(&id); err != nil {
		return CombatRecord{}, err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO combat_profiles(player_id) VALUES($1) ON CONFLICT DO NOTHING", attackerID); err != nil {
		return CombatRecord{}, err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO combat_profiles(player_id) VALUES($1) ON CONFLICT DO NOTHING", defenderID); err != nil {
		return CombatRecord{}, err
	}
	var level int
	var xp int64
	if err = tx.QueryRowContext(ctx, "SELECT level,xp FROM player_profiles WHERE id=$1 FOR UPDATE", attackerID).Scan(&level, &xp); err != nil {
		return CombatRecord{}, err
	}
	xp += result.XPReward
	for level < 100 {
		need := int64(level * level * 100)
		if xp < need {
			break
		}
		xp -= need
		level++
	}
	if _, err = tx.ExecContext(ctx, "UPDATE player_profiles SET level=$2,xp=$3,updated_at=NOW() WHERE id=$1", attackerID, level, xp); err != nil {
		return CombatRecord{}, err
	}
	if result.Winner == 1 {
		_, err = tx.ExecContext(ctx, "UPDATE combat_profiles SET wins=wins+1,rating=rating+10,updated_at=NOW() WHERE player_id=$1", attackerID)
		if err != nil {
			return CombatRecord{}, err
		}
		_, err = tx.ExecContext(ctx, "UPDATE combat_profiles SET losses=losses+1,rating=GREATEST(0,rating-10),updated_at=NOW() WHERE player_id=$1", defenderID)
	} else if result.Winner == 2 {
		_, err = tx.ExecContext(ctx, "UPDATE combat_profiles SET losses=losses+1,rating=GREATEST(0,rating-10),updated_at=NOW() WHERE player_id=$1", attackerID)
		if err != nil {
			return CombatRecord{}, err
		}
		_, err = tx.ExecContext(ctx, "UPDATE combat_profiles SET wins=wins+1,rating=rating+10,updated_at=NOW() WHERE player_id=$1", defenderID)
	}
	if err != nil {
		return CombatRecord{}, err
	}
	if err = tx.Commit(); err != nil {
		return CombatRecord{}, err
	}
	var w *string
	if result.Winner == 1 {
		w = &attackerID
	} else if result.Winner == 2 {
		w = &defenderID
	}
	return CombatRecord{ID: id, AttackerID: attackerID, DefenderID: defenderID, WinnerID: w, AttackerScore: result.AttackerScore, DefenderScore: result.DefenderScore, EnergyCost: result.EnergyCost, XPReward: result.XPReward, CreatedAt: time.Now().UTC()}, nil
}

func (r CombatRepository) History(ctx context.Context, playerID string) ([]CombatRecord, error) {
	rows, err := r.DB.SQL.QueryContext(ctx, "SELECT id::text,attacker_id::text,defender_id::text,winner_id::text,attacker_score,defender_score,energy_cost,xp_reward,created_at FROM combat_matches WHERE attacker_id=$1 OR defender_id=$1 ORDER BY created_at DESC LIMIT 50", playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []CombatRecord{}
	for rows.Next() {
		var x CombatRecord
		var w sql.NullString
		if err := rows.Scan(&x.ID, &x.AttackerID, &x.DefenderID, &w, &x.AttackerScore, &x.DefenderScore, &x.EnergyCost, &x.XPReward, &x.CreatedAt); err != nil {
			return nil, err
		}
		if w.Valid {
			x.WinnerID = &w.String
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
