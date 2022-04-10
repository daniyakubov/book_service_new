package cache

type Cache interface {
	Push(key string, value string) error
	Get(key string) ([]string, error)
}
