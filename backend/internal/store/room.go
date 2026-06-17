package store

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomStore struct {
	db *pgxpool.Pool
}

func NewRoomStore(db *pgxpool.Pool) *RoomStore {
	return &RoomStore{db: db}
}

//func (s *RoomStore) GetRoomsByBranch(ctx context.Context, branchID string)
