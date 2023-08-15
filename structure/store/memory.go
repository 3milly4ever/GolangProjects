package store

import (
	"database/sql"
	"structure/types"
)

type MySQLStore struct {
	db *sql.DB
}

// func NewMemoryStore() *MemoryStore {
// 	return &MemoryStore{}
// }

func (s *MySQLStore) Get(id int) *types.User {
	user := &types.User{}
	err := 

}