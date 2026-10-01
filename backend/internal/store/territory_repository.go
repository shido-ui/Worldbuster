package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrTerritoryNotFound = errors.New("territory not found")
var ErrInvalidInfluence = errors.New("invalid influence change")

type TerritoryRecord struct {
	ID                        string    `json:"id"`
	LocationID                string    `json:"locationId"`
	Name                      string    `json:"name"`
	ControllingOrganizationID *string   `json:"controllingOrganizationId,omitempty"`
	Influence                 int       `json:"influence"`
	Stability                 int       `json:"stability"`
	UpdatedAt                 time.Time `json:"updatedAt"`
}
type TerritoryInfluenceRecord struct {
	TerritoryID    string    `json:"territoryId"`
	OrganizationID string    `json:"organizationId"`
	Influence      int       `json:"influence"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
type TerritoryRepository struct{ DB *DB }

func (r TerritoryRepository) List(ctx context.Context) ([]TerritoryRecord, error) {
	rows, err := r.DB.SQL.QueryContext(ctx, "SELECT id::text,location_id,name,controlling_organization_id::text,influence,stability,updated_at FROM territories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	out := []TerritoryRecord{}
	for rows.Next() {
		var x TerritoryRecord
		var control sql.NullString
		if err := rows.Scan(&x.ID, &x.LocationID, &x.Name, &control, &x.Influence, &x.Stability, &x.UpdatedAt); err != nil {
			return nil, err
		}
		if control.Valid {
			x.ControllingOrganizationID = &control.String
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r TerritoryRepository) Get(ctx context.Context, id string) (TerritoryRecord, error) {
	var x TerritoryRecord
	var control sql.NullString
	err := r.DB.SQL.QueryRowContext(ctx, "SELECT id::text,location_id,name,controlling_organization_id::text,influence,stability,updated_at FROM territories WHERE id=$1", id).Scan(&x.ID, &x.LocationID, &x.Name, &control, &x.Influence, &x.Stability, &x.UpdatedAt)
	if err != nil {
		return x, err
	}
	if control.Valid {
		x.ControllingOrganizationID = &control.String
	}
	return x, nil
}

func (r TerritoryRepository) Influence(ctx context.Context, territoryID, orgID string) (int, error) {
	var n int
	err := r.DB.SQL.QueryRowContext(ctx, "SELECT influence FROM territory_influence WHERE territory_id=$1 AND organization_id=$2", territoryID, orgID).Scan(&n)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return n, err
}

func (r TerritoryRepository) AddInfluence(ctx context.Context, territoryID, orgID, actorID string, delta int, reason string) (TerritoryRecord, error) {
	if delta == 0 || delta > 100 || delta < -100 {
		return TerritoryRecord{}, ErrInvalidInfluence
	}
	tx, err := r.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return TerritoryRecord{}, err
	}
	defer func() { _ = tx.Rollback() }()
	var exists bool
	if err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM territories WHERE id=$1)", territoryID).Scan(&exists); err != nil || !exists {
		return TerritoryRecord{}, ErrTerritoryNotFound
	}
	var member bool
	if err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM organization_members WHERE organization_id=$1 AND character_id=$2)", orgID, actorID).Scan(&member); err != nil {
		return TerritoryRecord{}, err
	}
	if !member {
		return TerritoryRecord{}, errors.New("actor is not an organization member")
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO territory_influence(territory_id,organization_id,influence,last_activity_at) VALUES($1,$2,$3,NOW()) ON CONFLICT(territory_id,organization_id) DO UPDATE SET influence=GREATEST(0,LEAST(1000,territory_influence.influence+$3)),updated_at=NOW(),last_activity_at=NOW()", territoryID, orgID, delta)
	if err != nil {
		return TerritoryRecord{}, err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO territory_history(territory_id,organization_id,actor_id,delta,reason) VALUES($1,$2,$3,$4,$5)", territoryID, orgID, actorID, delta, reason)
	if err != nil {
		return TerritoryRecord{}, err
	}
	var x TerritoryRecord
	var control sql.NullString
	err = tx.QueryRowContext(ctx, "SELECT id::text,location_id,name,controlling_organization_id::text,influence,stability,updated_at FROM territories WHERE id=$1", territoryID).Scan(&x.ID, &x.LocationID, &x.Name, &control, &x.Influence, &x.Stability, &x.UpdatedAt)
	if err != nil {
		return TerritoryRecord{}, err
	}
	if control.Valid {
		x.ControllingOrganizationID = &control.String
	}
	if err = tx.Commit(); err != nil {
		return TerritoryRecord{}, err
	}
	return x, nil
}
