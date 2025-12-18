package policy

import (
	"path/filepath"
	"testing"
)

func TestRulesRepositoryLoadAndFilter(t *testing.T) {
	dir := filepath.Join("..", "..", "policies")
	repo, err := NewRulesRepository(dir)
	if err != nil {
		t.Fatalf("load rules: %v", err)
	}
	all := repo.List(map[string]string{})
	if len(all) == 0 {
		t.Fatalf("expected rules > 0")
	}
	eu := repo.List(map[string]string{"jurisdiction": "EU"})
	if len(eu) == 0 {
		t.Fatalf("expected EU rules > 0")
	}
}

