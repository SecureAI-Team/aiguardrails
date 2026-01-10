package rules

import (
	"errors"
	"sync"
)

// MemoryStore implements Store in memory.
type MemoryStore struct {
	mu    sync.RWMutex
	rules map[string]Rule
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rules: make(map[string]Rule),
	}
}

func (s *MemoryStore) Add(rule Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.rules[rule.ID]; exists {
		return errors.New("rule already exists")
	}
	s.rules[rule.ID] = rule
	return nil
}

func (s *MemoryStore) Get(id string) (*Rule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	rule, ok := s.rules[id]
	if !ok {
		return nil, errors.New("rule not found")
	}
	return &rule, nil
}

func (s *MemoryStore) List() ([]Rule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var list []Rule
	for _, r := range s.rules {
		list = append(list, r)
	}
	return list, nil
}

func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if rule, exists := s.rules[id]; exists && rule.IsSystem {
		return errors.New("cannot delete system rule")
	}
	delete(s.rules, id)
	return nil
}

func (s *MemoryStore) Update(rule Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.rules[rule.ID]; !exists {
		return errors.New("rule not found")
	}
	s.rules[rule.ID] = rule
	return nil
}
