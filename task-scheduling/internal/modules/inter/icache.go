package inter

type ICache interface {
	SetCache(key []byte, val []byte, expireSeconds int) error
	GetCache(key []byte) ([]byte, error)
	ClearCache()
}
