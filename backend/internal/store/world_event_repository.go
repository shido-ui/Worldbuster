package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type WorldEventDefinition struct {
	Code, Title, Description, EventType string
	CooldownSeconds                     int
	Active                              bool
}
type WorldEventInstance struct {
	ID, Code, Status string
	Severity         int
	LocationID       *string
	Payload          map[string]any
	StartedAt        time.Time
	ResolvedAt       *time.Time
}
type WorldEventRepository struct{ DB *DB }

func (r WorldEventRepository) Definitions(ctx context.Context) ([]WorldEventDefinition, error) {
	rows, err := r.DB.SQL.QueryContext(ctx, "SELECT code,title,description,event_type,cooldown_seconds,active FROM world_event_definitions WHERE active ORDER BY code")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []WorldEventDefinition{}
	for rows.Next() {
		var x WorldEventDefinition
		if err := rows.Scan(&x.Code, &x.Title, &x.Description, &x.EventType, &x.CooldownSeconds, &x.Active); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r WorldEventRepository) Create(ctx context.Context, code string, severity int, location string, payload map[string]any) (WorldEventInstance, error) {
	if severity < 1 {
		severity = 1
	}
	if severity > 5 {
		severity = 5
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return WorldEventInstance{}, err
	}
	var x WorldEventInstance
	var loc sql.NullString
	err = r.DB.SQL.QueryRowContext(ctx, "INSERT INTO world_event_instances(code,severity,location_id,payload) VALUES($1,$2,$3,$4) RETURNING id::text,code,status,severity,location_id,payload,started_at,resolved_at", code, severity, location, b).Scan(&x.ID, &x.Code, &x.Status, &x.Severity, &loc, &b, &x.StartedAt, &x.ResolvedAt)
	if loc.Valid {
		x.LocationID = &loc.String
	}
	if err != nil {
		return x, err
	}
	if err = json.Unmarshal(b, &x.Payload); err != nil {
		return x, err
	}
	return x, nil
}

func (r WorldEventRepository) Recent(ctx context.Context) ([]WorldEventInstance, error) {
	rows, err := r.DB.SQL.QueryContext(ctx, "SELECT id::text,code,status,severity,location_id,payload,started_at,resolved_at FROM world_event_instances ORDER BY started_at DESC LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []WorldEventInstance{}
	for rows.Next() {
		var x WorldEventInstance
		var loc sql.NullString
		var b []byte
		if err := rows.Scan(&x.ID, &x.Code, &x.Status, &x.Severity, &loc, &b, &x.StartedAt, &x.ResolvedAt); err != nil {
			return nil, err
		}
		if loc.Valid {
			x.LocationID = &loc.String
		}
		if err := json.Unmarshal(b, &x.Payload); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
