package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrOrganizationNotFound = errors.New("organization not found")
var ErrOrganizationFull = errors.New("organization is full")
var ErrOrganizationMember = errors.New("already a member")
var ErrOrganizationNotMember = errors.New("not a member")

type OrganizationRecord struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	OwnerID    string    `json:"ownerId"`
	MaxMembers int       `json:"maxMembers"`
	Level      int       `json:"level"`
	Reputation int       `json:"reputation"`
	Treasury   int64     `json:"treasury"`
	CreatedAt  time.Time `json:"createdAt"`
}
type OrganizationMemberRecord struct {
	OrganizationID string    `json:"organizationId"`
	PlayerID       string    `json:"playerId"`
	Role           string    `json:"role"`
	JoinedAt       time.Time `json:"joinedAt"`
}
type OrganizationRepository struct{ DB *DB }

func (r OrganizationRepository) List(ctx context.Context) ([]OrganizationRecord, error) {
	rows, err := r.DB.SQL.QueryContext(ctx, "SELECT id::text,name,type,owner_id::text,max_members,level,reputation,treasury,created_at FROM organizations ORDER BY reputation DESC,name LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []OrganizationRecord{}
	for rows.Next() {
		var x OrganizationRecord
		if err := rows.Scan(&x.ID, &x.Name, &x.Type, &x.OwnerID, &x.MaxMembers, &x.Level, &x.Reputation, &x.Treasury, &x.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r OrganizationRepository) Get(ctx context.Context, id string) (OrganizationRecord, error) {
	var x OrganizationRecord
	err := r.DB.SQL.QueryRowContext(ctx, "SELECT id::text,name,type,owner_id::text,max_members,level,reputation,treasury,created_at FROM organizations WHERE id=$1", id).Scan(&x.ID, &x.Name, &x.Type, &x.OwnerID, &x.MaxMembers, &x.Level, &x.Reputation, &x.Treasury, &x.CreatedAt)
	return x, err
}

func (r OrganizationRepository) Join(ctx context.Context, orgID, playerID, role string) (OrganizationMemberRecord, error) {
	if role != "" && role != "member" {
		return OrganizationMemberRecord{}, errors.New("invalid membership role")
	}
	role = "member"
	tx, err := r.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return OrganizationMemberRecord{}, err
	}
	defer tx.Rollback()
	var max, count int
	if err = tx.QueryRowContext(ctx, "SELECT max_members,(SELECT COUNT(*) FROM organization_members WHERE organization_id=organizations.id) FROM organizations WHERE id=$1 FOR UPDATE", orgID).Scan(&max, &count); err != nil {
		return OrganizationMemberRecord{}, ErrOrganizationNotFound
	}
	if count >= max {
		return OrganizationMemberRecord{}, ErrOrganizationFull
	}
	var x OrganizationMemberRecord
	err = tx.QueryRowContext(ctx, "INSERT INTO organization_members(organization_id,character_id,role) VALUES($1,$2,$3) RETURNING organization_id::text,character_id::text,role,joined_at", orgID, playerID, role).Scan(&x.OrganizationID, &x.PlayerID, &x.Role, &x.JoinedAt)
	if err != nil {
		return OrganizationMemberRecord{}, ErrOrganizationMember
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO organization_history(organization_id,actor_id,event_type,payload) VALUES($1,$2,'member.joined','{}')", orgID, playerID); err != nil {
		return OrganizationMemberRecord{}, err
	}
	if err = tx.Commit(); err != nil {
		return OrganizationMemberRecord{}, err
	}
	return x, nil
}

func (r OrganizationRepository) Membership(ctx context.Context, orgID, playerID string) (OrganizationMemberRecord, error) {
	var x OrganizationMemberRecord
	err := r.DB.SQL.QueryRowContext(ctx, "SELECT organization_id::text,character_id::text,role,joined_at FROM organization_members WHERE organization_id=$1 AND character_id=$2", orgID, playerID).Scan(&x.OrganizationID, &x.PlayerID, &x.Role, &x.JoinedAt)
	return x, err
}

func (r OrganizationRepository) FactionReputation(ctx context.Context, orgID, playerID string) (int, error) {
	var n int
	err := r.DB.SQL.QueryRowContext(ctx, "SELECT score FROM organization_reputation WHERE organization_id=$1 AND player_id=$2", orgID, playerID).Scan(&n)
	if err == sql.ErrNoRows {
		_, err = r.DB.SQL.ExecContext(ctx, "INSERT INTO organization_reputation(organization_id,player_id) VALUES($1,$2) ON CONFLICT DO NOTHING", orgID, playerID)
		if err != nil {
			return 0, err
		}
		return 0, nil
	}
	return n, err
}

func (r OrganizationRepository) ChangeFactionReputation(ctx context.Context, orgID, playerID string, delta int, reason string) (int, error) {
	if delta > 100 || delta < -100 {
		return 0, errors.New("reputation delta too large")
	}
	tx, err := r.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO organization_reputation(organization_id,player_id) VALUES($1,$2) ON CONFLICT DO NOTHING", orgID, playerID)
	if err != nil {
		return 0, err
	}
	var score int
	err = tx.QueryRowContext(ctx, "UPDATE organization_reputation SET score=GREATEST(-1000,LEAST(1000,score+$3)),updated_at=NOW() WHERE organization_id=$1 AND player_id=$2 RETURNING score", orgID, playerID, delta).Scan(&score)
	if err != nil {
		return 0, err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO organization_history(organization_id,actor_id,event_type,payload) VALUES($1,$2,'faction.reputation.changed',jsonb_build_object('delta',$3,'reason',$4))", orgID, playerID, delta, reason)
	if err != nil {
		return 0, err
	}
	err = tx.Commit()
	return score, err
}

func (r OrganizationRepository) IsMember(ctx context.Context, orgID, playerID string) (bool, error) {
	var ok bool
	err := r.DB.SQL.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM organization_members WHERE organization_id=$1 AND character_id=$2)", orgID, playerID).Scan(&ok)
	return ok, err
}

func (r OrganizationRepository) UpdateActivity(ctx context.Context, orgID string, delta int) (int, error) {
	if delta > 100 || delta < -100 {
		return 0, errors.New("activity delta too large")
	}
	var score int
	err := r.DB.SQL.QueryRowContext(ctx, "UPDATE organizations SET activity_score=GREATEST(0,LEAST(1000,activity_score+$2)),influence=GREATEST(0,influence+$2) WHERE id=$1 RETURNING activity_score", orgID, delta).Scan(&score)
	if err == sql.ErrNoRows {
		return 0, ErrOrganizationNotFound
	}
	return score, err
}

func (r OrganizationRepository) ProposeAlliance(ctx context.Context, orgID, targetID, actorID string) error {
	if orgID == targetID {
		return errors.New("organization cannot ally with itself")
	}
	tx, err := r.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var role string
	if err = tx.QueryRowContext(ctx, "SELECT role FROM organization_members WHERE organization_id=$1 AND character_id=$2", orgID, actorID).Scan(&role); err != nil {
		return ErrOrganizationNotMember
	}
	if role != "owner" && role != "leader" {
		return errors.New("insufficient organization role")
	}
	var exists bool
	if err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM organizations WHERE id=$1)", targetID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return ErrOrganizationNotFound
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO organization_alliances(organization_id,target_organization_id,status) VALUES($1,$2,'PROPOSED') ON CONFLICT(organization_id,target_organization_id) DO UPDATE SET status='PROPOSED',updated_at=NOW()", orgID, targetID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO organization_history(organization_id,actor_id,event_type,payload) VALUES($1,$2,'alliance.proposed',jsonb_build_object('targetOrganizationId',$3))", orgID, actorID, targetID); err != nil {
		return err
	}
	return tx.Commit()
}

func (r OrganizationRepository) SetAllianceStatus(ctx context.Context, orgID, targetID, actorID, status string) error {
	if status != "ACTIVE" && status != "ENDED" {
		return errors.New("invalid alliance status")
	}
	tx, err := r.DB.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var role string
	if err = tx.QueryRowContext(ctx, "SELECT role FROM organization_members WHERE organization_id=$1 AND character_id=$2", orgID, actorID).Scan(&role); err != nil {
		return ErrOrganizationNotMember
	}
	if role != "owner" && role != "leader" {
		return errors.New("insufficient organization role")
	}
	result, err := tx.ExecContext(ctx, "UPDATE organization_alliances SET status=$3,updated_at=NOW() WHERE organization_id=$1 AND target_organization_id=$2", orgID, targetID, status)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return errors.New("alliance not found")
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO organization_history(organization_id,actor_id,event_type,payload) VALUES($1,$2,'alliance.status_changed',jsonb_build_object('targetOrganizationId',$3,'status',$4))", orgID, actorID, targetID, status)
	if err != nil {
		return err
	}
	return tx.Commit()
}
