package storage

const (
	action_null                             = ""
	action_cache_hmget                      = "cache.HMGet"
	action_cache_hmset                      = "cache.HMSet"
	action_cache_hdel                       = "cache.HDel"
	action_cache_del                        = "cache.Del"
	action_cache_zscore                     = "cache.ZScore"
	action_cache_zrank                      = "cache.ZRank"
	action_cache_zrange                     = "cache.ZRange"
	action_cache_zrevrange                  = "cache.ZRevRange"
	action_cache_zrangebyscore              = "cache.ZRangeByScore"
	action_cache_zrevrangebyscore           = "cache.ZRevRangeByScore"
	action_cache_zrangebylex                = "cache.ZRangeByLex"
	action_cache_zrevrangebylex             = "cache.ZRevRangeByLex"
	action_cache_zrangewithscores           = "cache.ZRangeWithScores"
	action_cache_zrevrangewithscores        = "cache.ZRevRangeWithScores"
	action_cache_zrangebyscorewithscores    = "cache.ZRangeByScoreWithScores"
	action_cache_zrevrangebyscorewithscores = "cache.ZRevRangeByScoreWithScores"
	action_db_insert                        = "db.Insert"
	action_db_update                        = "db.Update"
	action_db_remove                        = "db.Remove"
	action_db_get                           = "db.Get"
)

func action_get_field(table, field string) string {
	return "table `" + table + "` GetField `" + field + "`"
}
func action_set_field(table, field string) string {
	return "table `" + table + "` SetField `" + field + "`"
}

type ErrorHandler func(action string, err error) error

type Engine interface {
	Repository
	// Cache returns CacheProxy
	Cache() CacheProxy
	// Database returns DatabaseProxy
	Database() DatabaseProxy
	// SetErrorHandler sets handler for handling error
	SetErrorHandler(eh ErrorHandler)
	// AddIndex adds an index
	AddIndex(index Index)
	// NewSession new a session
	NewSession() Session
}

func JoinKey(engineName, originKey string) string {
	return engineName + "@" + originKey
}

func JoinField(key, originField string) string {
	return key + ":" + originField
}

func JoinIndexKey(engineName string, index Index) string {
	return engineName + "@" + index.Table() + ":" + index.Name()
}

type GetOption func(*getOptions)

type getOptions struct {
	syncFromDatabase bool
}

func WithSyncFromDatabase() GetOption {
	return syncFromDatabase
}

func syncFromDatabase(opt *getOptions) {
	opt.syncFromDatabase = true
}

type RangeOption func(*rangeOptions)

type rangeOptions struct {
	withScores bool
	rev        bool
	offset     int64
	count      int64
}

const (
	rangeByScore = 0
	rangeByLex   = 1
)

func RangeRev() RangeOption {
	return func(opts *rangeOptions) {
		opts.rev = true
	}
}

func RangeWithScores() RangeOption {
	return func(opts *rangeOptions) {
		opts.withScores = true
	}
}

func RangeOffset(offset int64) RangeOption {
	return func(opts *rangeOptions) {
		opts.offset = offset
	}
}

func RangeCount(count int64) RangeOption {
	return func(opts *rangeOptions) {
		opts.count = count
	}
}

type RangeLexResult interface {
	KeyList
}

type RangeResult interface {
	KeyList
	Score(i int) float64
}

// engine implements Engine interface
type engine struct {
	name         string
	cache        CacheProxy
	database     DatabaseProxy
	errorHandler ErrorHandler
	indexes      map[string]map[string]Index
}

// NewEngine creates an engine which named name.
// Parameter database MUST be not nil, but cache can be nil.
func NewEngine(name string, database DatabaseProxy, cache CacheProxy) Engine {
	eng := &engine{
		name:     name,
		database: database,
		cache:    cache,
		indexes:  make(map[string]map[string]Index),
	}
	return eng
}

func (eng *engine) Name() string            { return eng.name }
func (eng *engine) Cache() CacheProxy       { return eng.cache }
func (eng *engine) Database() DatabaseProxy { return eng.database }

// SetErrorHandler sets handler for handling error
func (eng *engine) SetErrorHandler(eh ErrorHandler) {
	eng.errorHandler = eh
}

// AddIndex adds an index
func (eng *engine) AddIndex(index Index) {
	tableName := index.Table()
	idx, ok := eng.indexes[tableName]
	if !ok {
		idx = make(map[string]Index)
		eng.indexes[tableName] = idx
	}
	idx[index.Name()] = index
}

func (eng *engine) NewSession() Session {
	return eng.newSession()
}

func (eng *engine) newSession() *session {
	s := &session{
		eng:      eng,
		database: eng.database.NewSession(),
	}
	if eng.cache != nil {
		s.cache = eng.cache.NewSession()
	}
	return s
}

// Insert inserts new records or updates all fields of records
func (eng *engine) Insert(tables ...Table) error {
	s := eng.newSession()
	defer s.Close()
	for _, table := range tables {
		action, err := s.update(table, true)
		if err != nil {
			return s.catch("Insert: "+action, err)
		}
	}
	return nil
}

// Update updates specific fields of record
func (eng *engine) Update(table Table, fields ...string) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.update(table, false, fields...)
	if err != nil {
		return s.catch("Update: "+action, err)
	}
	return nil
}

// Find gets many records
func (eng *engine) Find(meta TableMeta, keys KeyList, setters FieldSetterList, fields ...string) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.find(meta, keys, setters, fields)
	if err != nil {
		return s.catch("Find: "+action, err)
	}
	return nil
}

// Get gets one record all fields
func (eng *engine) Get(table Table, opts ...GetOption) (bool, error) {
	s := eng.newSession()
	defer s.Close()
	opt := getOptions{}
	for _, o := range opts {
		o(&opt)
	}
	action, ok, err := s.get(table, opt, table.Meta().Fields()...)
	if err != nil {
		return ok, s.catch("Get: "+action, err)
	}
	return ok, nil
}

// Get gets one record specific fields
func (eng *engine) GetFields(table Table, fields ...string) (bool, error) {
	s := eng.newSession()
	defer s.Close()
	opt := getOptions{}
	action, ok, err := s.get(table, opt, fields...)
	if err != nil {
		return ok, s.catch("Get: "+action, err)
	}
	return ok, nil
}

// Remove removes one record
func (eng *engine) Remove(table ReadonlyTable) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.remove(table.Meta(), table.Key())
	if err != nil {
		return s.catch("Remove: "+action, err)
	}
	return nil
}

// RemoveRecords removes records by keys
func (eng *engine) RemoveRecords(meta TableMeta, keys ...interface{}) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.remove(meta, keys...)
	if err != nil {
		return s.catch("RemoveRecords: "+action, err)
	}
	return nil
}

// Clear removes all records of table
func (eng *engine) Clear(table string) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.clear(table)
	if err != nil {
		return s.catch("Clear "+table+": "+action, err)
	}
	return nil
}

// FindView loads view by keys and store loaded data to setters
func (eng *engine) FindView(view View, keys KeyList, setters FieldSetterList) error {
	s := eng.newSession()
	defer s.Close()
	action, err := s.recursivelyLoadView(view, keys, setters)
	if err != nil {
		return s.catch("FindView: "+action, err)
	}
	return nil
}

// IndexRank gets rank of table key in index, returns InvalidRank if key not found
func (eng *engine) IndexRank(index Index, key interface{}) (int64, error) {
	s := eng.newSession()
	defer s.Close()
	action, rank, err := s.indexRank(index, key)
	if err != nil {
		return rank, s.catch("IndexRank: "+action, err)
	}
	return rank, nil
}

// IndexScore gets score of table key in index, returns InvalidScore if key not found
func (eng *engine) IndexScore(index Index, key interface{}) (int64, error) {
	s := eng.newSession()
	defer s.Close()
	action, score, err := s.indexScore(index, key)
	if err != nil {
		return score, s.catch("IndexScore: "+action, err)
	}
	return score, nil
}

func (eng *engine) IndexRange(index Index, start, stop int64, opts ...RangeOption) (RangeResult, error) {
	s := eng.newSession()
	defer s.Close()
	action, result, err := s.indexRange(index, start, stop, s.applyRangeOption(opts))
	if err != nil {
		return nil, s.catch("IndexRange: "+action, err)
	}
	return result, nil
}

func (eng *engine) IndexRangeByScore(index Index, min, max float64, opts ...RangeOption) (RangeResult, error) {
	s := eng.newSession()
	defer s.Close()
	action, result, err := s.indexRangeByScore(index, min, max, s.applyRangeOption(opts))
	if err != nil {
		return nil, s.catch("IndexRangeByScore: "+action, err)
	}
	return result, nil
}

func (eng *engine) IndexRangeByLex(index Index, min, max string, opts ...RangeOption) (RangeLexResult, error) {
	s := eng.newSession()
	defer s.Close()
	action, result, err := s.indexRangeByLex(index, min, max, s.applyRangeOption(opts))
	if err != nil {
		return nil, s.catch("IndexRangeByLex: "+action, err)
	}
	return result, nil
}
