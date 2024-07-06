package services

type CacheService interface {
	Put(key string, value string) error
	Get(key string) (string, error)
}
