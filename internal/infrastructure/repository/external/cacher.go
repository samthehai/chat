package external

type Cacher interface {
	LPush(key string, values []byte) error
	SAdd(key string, values []byte) error
	LRange(key string, start, stop int64) ([]string, error)
	SMembers(key string) ([]string, error)
}
