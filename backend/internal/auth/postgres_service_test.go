package auth

import "testing"

func TestPostgresServiceShape(t *testing.T){
 var _ *PostgresService = NewPostgresService(nil,nil)
}
