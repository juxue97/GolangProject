package main

import "context"

type store struct {
	// MongoDB
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context) error {
	return nil
}
