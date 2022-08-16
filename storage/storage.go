package storage

import "errors"

type (
	Storage interface {
		Store(orderUID, jsonOrder string) error
	}

	OrderDB struct {
		order_uid string
		data      string
	}
)

var (
	ErrAlreadyExists = errors.New("order already exists.")
)
