package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrJobNotFound = errors.New("job not found")
var ErrAlreadyEmployed = errors.New("already employed")
var ErrJobRequirement = errors.New("job requirements not met")

type JobRepository struct{ DB *DB }
type JobRecord struct {
	ID                string
	Name              string
	Department        string
	BaseSalary        int64
	RequiredLevel     int
	RequiredStat      int
	RequiredEducation int
}

func (r JobRepository) List(ctx context.Context) ([]JobRecord, error) {
	rows, err := r.DB.SQL.QueryContext(ctx, "SELECT id::text,name,department,base_salary,required_level,required_stat,required_education FROM jobs ORDER BY required_level, name")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []JobRecord{}
	for rows.Next() {
		var j JobRecord
		if err := rows.Scan(&j.ID, &j.Name, &j.Department, &j.BaseSalary, &j.RequiredLevel, &j.RequiredStat, &j.RequiredEducation); err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	return out, rows.Err()
}
func (r JobRepository) Employment(ctx context.Context, playerID string) (sql.NullString, error) {
	var id sql.NullString
	err := r.DB.SQL.QueryRowContext(ctx, "SELECT job_id::text FROM player_employment WHERE player_id=$1", playerID).Scan(&id)
	return id, err
}
func (r JobRepository) Employ(ctx context.Context, playerID, jobID string, level, stat, education int) (JobRecord, error) {
	var j JobRecord
	err := r.DB.WithTx(ctx, func(tx *sql.Tx) error {
		var existing string
		err := tx.QueryRowContext(ctx, "SELECT job_id::text FROM player_employment WHERE player_id=$1 FOR UPDATE", playerID).Scan(&existing)
		if err == nil {
			return ErrAlreadyEmployed
		}
		if err != sql.ErrNoRows {
			return err
		}
		err = tx.QueryRowContext(ctx, "SELECT id::text,name,department,base_salary,required_level,required_stat,required_education FROM jobs WHERE id=$1", jobID).Scan(&j.ID, &j.Name, &j.Department, &j.BaseSalary, &j.RequiredLevel, &j.RequiredStat, &j.RequiredEducation)
		if err == sql.ErrNoRows {
			return ErrJobNotFound
		}
		if err != nil {
			return err
		}
		if level < j.RequiredLevel || stat < j.RequiredStat || education < j.RequiredEducation {
			return ErrJobRequirement
		}
		_, err = tx.ExecContext(ctx, "INSERT INTO player_employment(player_id,job_id) VALUES($1,$2)", playerID, jobID)
		return err
	})
	return j, err
}

type SalaryPayment struct {
	PlayerID  string
	AccountID string
	JobID     string
	Amount    int64
	Level     int
	XP        int64
}

func (r JobRepository) SettleDueSalaries(ctx context.Context, intervalMinutes int) ([]SalaryPayment, error) {
	if intervalMinutes <= 0 {
		intervalMinutes = 60
	}
	payments := []SalaryPayment{}
	err := r.DB.WithTx(ctx, func(tx *sql.Tx) error {
		rows, err := tx.QueryContext(ctx, fmt.Sprintf(`SELECT pe.player_id::text,p.account_id::text,pe.job_id::text,j.base_salary
   FROM player_employment pe
   JOIN player_profiles p ON p.id=pe.player_id
   JOIN jobs j ON j.id=pe.job_id
   WHERE pe.last_paid_at <= NOW() - INTERVAL '%d minutes'
   ORDER BY pe.last_paid_at
   FOR UPDATE OF pe SKIP LOCKED LIMIT 100`, intervalMinutes))
		if err != nil {
			return err
		}
		defer func() { _ = rows.Close() }()
		type due struct {
			player, account, job string
			salary               int64
		}
		var ds []due
		for rows.Next() {
			var d due
			if err := rows.Scan(&d.player, &d.account, &d.job, &d.salary); err != nil {
				return err
			}
			ds = append(ds, d)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		for _, d := range ds {
			var balance int64
			var accountID string
			if err := tx.QueryRowContext(ctx, `SELECT id::text,balance FROM economy_accounts WHERE owner_account_id=$1 FOR UPDATE`, d.account).Scan(&accountID, &balance); err != nil {
				return err
			}
			balance += d.salary
			if _, err := tx.ExecContext(ctx, `UPDATE economy_accounts SET balance=$1,updated_at=NOW() WHERE id=$2`, balance, accountID); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO economy_ledger(id,account_id,amount,balance_after,reason,reference_id) VALUES(gen_random_uuid(),$1,$2,$3,'job_salary',$4)`, accountID, d.salary, balance, d.player); err != nil {
				return err
			}
			var level int
			var xp int64
			if err := tx.QueryRowContext(ctx, `SELECT level,xp FROM player_profiles WHERE id=$1 FOR UPDATE`, d.player).Scan(&level, &xp); err != nil {
				return err
			}
			xp += 25
			for level < 100 {
				need := int64(level * level * 100)
				if xp < need {
					break
				}
				xp -= need
				level++
			}
			if _, err := tx.ExecContext(ctx, `UPDATE player_profiles SET level=$2,xp=$3,updated_at=NOW() WHERE id=$1`, d.player, level, xp); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `UPDATE player_employment SET last_paid_at=NOW() WHERE player_id=$1`, d.player); err != nil {
				return err
			}
			payments = append(payments, SalaryPayment{PlayerID: d.player, AccountID: d.account, JobID: d.job, Amount: d.salary, Level: level, XP: xp})
		}
		return nil
	})
	return payments, err
}
