package goredisproxy

import (
	"gopkg.in/redis.v5"

	"github.com/midlang/mid/x/go/storage"
	"github.com/mkideal/pkg/typeconv"
)

// proxy implements storage.CacheProxy
type proxy struct {
	client redis.Cmdable
}

func New(client redis.Cmdable) storage.CacheProxy {
	return &proxy{client: client}
}

func (p *proxy) NewSession() storage.CacheProxySession {
	return &proxySession{client: p.client}
}

// proxySession implements storage.CacheProxySession
type proxySession struct {
	client redis.Cmdable
}

func (p *proxySession) Begin() error    { return nil }
func (p *proxySession) Commit() error   { return nil }
func (p *proxySession) Rollback() error { return nil }
func (p *proxySession) Close()          {}

func (p *proxySession) error(err error) error {
	if err == redis.Nil {
		return storage.ErrNotFound
	}
	return err
}

func (p *proxySession) TableName(engineName, table string) string {
	return engineName + "@" + table
}

func (p *proxySession) FieldName(engineName, table string, key interface{}, field string) string {
	return typeconv.ToString(key) + ":" + field
}

func (p *proxySession) IndexKey(engineName string, index storage.Index) string {
	return engineName + "@" + index.Table() + ":" + index.Name()
}

func (p *proxySession) HDel(table string, fields ...string) (int64, error) {
	x, err := p.client.HDel(table, fields...).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) HExists(table, field string) (bool, error) {
	x, err := p.client.HExists(table, field).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) HIncrBy(table, field string, incr int64) (int64, error) {
	x, err := p.client.HIncrBy(table, field, incr).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) HMGet(table string, fields ...string) ([]interface{}, error) {
	x, err := p.client.HMGet(table, fields...).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) HMSet(table string, fields map[string]string) (string, error) {
	x, err := p.client.HMSet(table, fields).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) Delete(keys ...string) (int64, error) {
	x, err := p.client.Del(keys...).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) ZAdd(key string, members ...redis.Z) (int64, error) {
	x, err := p.client.ZAdd(key, members...).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) ZRem(key string, members ...interface{}) (int64, error) {
	x, err := p.client.ZRem(key, members...).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) ZRank(key, member string) (int64, error) {
	x, err := p.client.ZRank(key, member).Result()
	err = p.error(err)
	return x, err
}

func (p *proxySession) ZScore(key, member string) (int64, error) {
	x, err := p.client.ZScore(key, member).Result()
	err = p.error(err)
	return int64(x), err
}
