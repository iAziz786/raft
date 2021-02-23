package storage

import "errors"

var (
	ErrNotFound = errors.New("key not found")
)

type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}
