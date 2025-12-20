package policy

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestRuleStoreAttachAndList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	store := NewRuleStore(db)

	mock.ExpectExec("INSERT INTO policy_rules").
		WithArgs("p1", "r1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := store.Attach("p1", "r1"); err != nil {
		t.Fatalf("attach failed: %v", err)
	}

	mock.ExpectQuery("SELECT rule_id FROM policy_rules").
		WithArgs("p1").
		WillReturnRows(sqlmock.NewRows([]string{"rule_id"}).AddRow("r1"))

	ids, err := store.ListByPolicy("p1")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(ids) != 1 || ids[0] != "r1" {
		t.Fatalf("unexpected ids: %v", ids)
	}
}

