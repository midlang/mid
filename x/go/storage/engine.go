package storage

const (
	action_null            = ""
	action_cache_hmget     = "cache.HMGET"
	action_cache_hmset     = "cache.HMSET"
	action_cache_hdel      = "cache.HDEL"
	action_cache_del       = "cache.DEL"
	action_cache_zscore    = "cache.ZSCORE"
	action_cache_zrank     = "cache.ZRANK"
	action_cache_zrange    = "cache.ZRANGE"
	action_cache_zrevrange = "cache.ZREVRANGE"
	action_db_insert       = "db.INSERT"
	action_db_update       = "db.UPDATE"
	action_db_remove       = "db.REMOVE"
	action_db_get          = "db.GET"
)

func action_get_field(table, field string) string {
	return "table `" + table + "` GetField `" + field + "`"
}
func action_set_field(table, field string) string {
	return "table `" + table + "` SetField `" + field + "`"
}

type ErrorHandler func(action string, err error) error

type Engine interface {
	OpAPI
	// Name returns database name
	Name() string
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

type GetOption func(*getOptions)

type getOptions struct {
	SyncFromDatabase bool
}

func WithSyncFromDatabase() GetOption {
	return syncFromDatabase
}

func syncFromDatabase(opt *getOptions) {
	opt.SyncFromDatabase = true
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
