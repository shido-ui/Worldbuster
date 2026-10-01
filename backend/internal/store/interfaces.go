package store

import (
	"context"
	"database/sql"
)

type JobStore interface {
	List(context.Context) ([]JobRecord, error)
	Employment(context.Context, string) (sql.NullString, error)
	Employ(context.Context, string, string, int, int, int) (JobRecord, error)
	SettleDueSalaries(context.Context, int) ([]SalaryPayment, error)
}

var _ JobStore = JobRepository{}
