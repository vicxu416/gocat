package server

type Cacher interface {
	Get(key []byte) ([]byte, error)
	Set(key, val []byte) error
	Del(key []byte) error
}
