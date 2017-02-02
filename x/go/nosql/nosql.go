package nosql

type Proxy interface {
	Get(key string) (string, error)
	Set(key, value string) error
	Hmget()
	Hmset()
}
