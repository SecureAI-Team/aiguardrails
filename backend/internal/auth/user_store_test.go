package auth

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"aiguardrails/internal/rbac"
)

func TestUserStoreCreateAndVerify(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").
		WillReturnResult(sqlmock.NewResult(1, 1))

	store := NewUserStore(db)
	user, err := store.Create("admin@example.com", "Passw0rd!", rbac.RolePlatformAdmin)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if user.Role != rbac.RolePlatformAdmin {
		t.Fatalf("expected role admin")
	}

	// prepare verify
	mock.ExpectQuery("SELECT id, username, password_hash, role, created_at, updated_at FROM users").
		WithArgs("admin@example.com").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "username", "password_hash", "role", "created_at", "updated_at"}).
				AddRow(user.ID, user.Username, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt),
		)

	out, err := store.Verify("admin@example.com", "Passw0rd!")
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if out.Username != user.Username {
		t.Fatalf("verify returned wrong user")
	}
}

func TestEnsureBootUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewUserStore(db)

	mock.ExpectExec("INSERT INTO users").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = store.EnsureBootUser("boot@example.com", "StrongPass123", rbac.RolePlatformAdmin)
	if err != nil {
		t.Fatalf("ensure boot user failed: %v", err)
	}
}

