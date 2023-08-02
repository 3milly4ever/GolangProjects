package store

import "structure/types"

type Store interface {
	Get(int) *types.User
}

// func (s *Store) Get(id int) *types.User {
// 	return &types.User{
// 		ID:   1,
// 		Name: "Foo",
// 	}
// }
