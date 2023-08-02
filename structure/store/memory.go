package store

import "structure/types"

type MemoryStore struct {
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Get(id int) *types.User {
	return &types.User{
		ID:   1,
		Name: "Foo",
	}
}
