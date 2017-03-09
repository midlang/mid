package storage

import (
	"github.com/mkideal/pkg/typeconv"
)

type Tx interface {
	Begin() error
	Commit() error
	Rollback() error
	Close()
}

type OpAPI interface {
	// Insert inserts new records
	Insert(tables ...Table) error
	// Update updates specific fields of record
	Update(table Table, fields ...string) error
	// Find gets records by keys
	Find(meta TableMeta, keys KeyList, setters FieldSetterList, fields ...string) error
	// Get gets one record all fields
	Get(table Table, opts ...GetOption) (bool, error)
	// Get gets one record specific fields
	GetFields(table Table, fields ...string) (bool, error)
	// Remove removes one record
	Remove(table ReadonlyTable) error
	// RemoveRecords removes records by keys
	RemoveRecords(meta TableMeta, keys ...interface{}) error
	// Clear removes all records of table
	Clear(table string) error
	// FindView loads view by keys and store loaded data to setters
	FindView(view View, keys KeyList, setters FieldSetterList) error
	// IndexRank gets rank of table key in index, returns InvalidRank if key not found
	IndexRank(index Index, key interface{}) (int64, error)
	// IndexScore gets score of table key in index, returns InvalidScore if key not found
	IndexScore(index Index, key interface{}) (int64, error)
}

type Session interface {
	Name() string
	Cache() CacheProxySession
	Database() DatabaseProxySession
	Tx
	OpAPI
}

type session struct {
	eng      *engine
	cache    CacheProxySession
	database DatabaseProxySession
}

func (s *session) catch(action string, err error) error {
	if err != nil && s.eng.errorHandler != nil {
		err = s.eng.errorHandler(action, err)
	}
	return err
}

func (s *session) Name() string                   { return s.eng.name }
func (s *session) Cache() CacheProxySession       { return s.cache }
func (s *session) Database() DatabaseProxySession { return s.database }

func (s *session) Begin() error {
	err := s.database.Begin()
	if err == nil && s.cache != nil {
		err = s.cache.Begin()
	}
	return err
}

func (s *session) Commit() error {
	err := s.database.Commit()
	if err == nil && s.cache != nil {
		err = s.cache.Commit()
	}
	return err
}

func (s *session) Rollback() error {
	err := s.database.Rollback()
	if err == nil && s.cache != nil {
		err = s.cache.Rollback()
	}
	return err
}

func (s *session) Close() {
	s.database.Close()
	if s.cache != nil {
		s.cache.Close()
	}
}

// Insert inserts new records or updates all fields of records
func (s *session) Insert(tables ...Table) error {
	for _, table := range tables {
		action, err := s.update(table, true)
		if err != nil {
			return s.catch("Insert: "+action, err)
		}
	}
	return nil
}

// Update updates specific fields of record
func (s *session) Update(table Table, fields ...string) error {
	action, err := s.update(table, false, fields...)
	if err != nil {
		return s.catch("Update: "+action, err)
	}
	return nil
}

// Find gets many records
func (s *session) Find(meta TableMeta, keys KeyList, setters FieldSetterList, fields ...string) error {
	action, err := s.find(meta, keys, setters, fields)
	if err != nil {
		return s.catch("Find: "+action, err)
	}
	return nil
}

// Get gets one record all fields
func (s *session) Get(table Table, opts ...GetOption) (bool, error) {
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
func (s *session) GetFields(table Table, fields ...string) (bool, error) {
	opt := getOptions{}
	action, ok, err := s.get(table, opt, fields...)
	if err != nil {
		return ok, s.catch("Get: "+action, err)
	}
	return ok, nil
}

// Remove removes one record
func (s *session) Remove(table ReadonlyTable) error {
	action, err := s.remove(table.Meta(), table.Key())
	if err != nil {
		return s.catch("Remove: "+action, err)
	}
	return nil
}

// RemoveRecords removes records by keys
func (s *session) RemoveRecords(meta TableMeta, keys ...interface{}) error {
	action, err := s.remove(meta, keys...)
	if err != nil {
		return s.catch("RemoveRecords: "+action, err)
	}
	return nil
}

// Clear removes all records of table
func (s *session) Clear(table string) error {
	action, err := s.clear(table)
	if err != nil {
		return s.catch("Clear "+table+": "+action, err)
	}
	return nil
}

// FindView loads view by keys and store loaded data to setters
func (s *session) FindView(view View, keys KeyList, setters FieldSetterList) error {
	action, err := s.recursivelyLoadView(view, keys, setters)
	if err != nil {
		return s.catch("FindView: "+action, err)
	}
	return nil
}

// IndexRank gets rank of table key in index, returns InvalidRank if key not found
func (s *session) IndexRank(index Index, key interface{}) (int64, error) {
	action, rank, err := s.indexRank(index, key)
	if err != nil {
		return rank, s.catch("IndexRank: "+action, err)
	}
	return rank, nil
}

// IndexScore gets score of table key in index, returns InvalidScore if key not found
func (s *session) IndexScore(index Index, key interface{}) (int64, error) {
	action, score, err := s.indexScore(index, key)
	if err != nil {
		return score, s.catch("IndexScore: "+action, err)
	}
	return score, nil
}

//----------------
// implementation
//----------------

func (s *session) update(table Table, insert bool, fields ...string) (string, error) {

	// database op
	if insert {
		_, err := s.database.Insert(table)
		if err != nil {
			return action_db_insert, err
		}
	} else {
		if len(fields) == 0 {
			fields = table.Meta().Fields()
		}
		_, err := s.database.Update(table, fields...)
		if err != nil {
			return action_db_update, err
		}
	}

	// cache op
	if s.cache == nil {
		return action_null, nil
	}
	return s.updateCache(table, fields...)
}

func (s *session) updateCache(table Table, fields ...string) (string, error) {
	var (
		meta = table.Meta()
		key  = table.Key()
	)
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	args := make(map[string]string)
	fieldKey := typeconv.ToString(key)
	for _, field := range fields {
		key := s.cache.FieldName(s.eng.name, meta.Name(), fieldKey, field)
		value, ok := table.GetField(field)
		if !ok {
			return action_get_field(meta.Name(), field), ErrFieldNotFound
		}
		args[key] = typeconv.ToString(value)
	}
	action, err := s.updateIndex(table, key, fields)
	if err != nil {
		return action, err
	}
	_, err = s.cache.HMSet(s.cache.TableName(s.eng.name, meta.Name()), args)
	return action_cache_hmset, err
}

func (s *session) remove(meta TableMeta, keys ...interface{}) (string, error) {
	// database op
	if _, err := s.database.Remove(meta.Name(), meta.Key(), keys...); err != nil {
		return action_db_update, err
	}

	// cache op
	fields := meta.Fields()
	if s.cache == nil {
		return action_null, nil
	}
	args := make([]string, 0, len(fields)*len(keys))
	for _, key := range keys {
		fieldKey := typeconv.ToString(key)
		for _, field := range fields {
			args = append(args, s.cache.FieldName(s.eng.name, meta.Name(), fieldKey, field))
		}
	}
	if action, err := s.removeIndex(meta.Name(), keys...); err != nil {
		return action, err
	}
	_, err := s.cache.HDel(s.cache.TableName(s.eng.name, meta.Name()), args...)
	return action_cache_hdel, err
}

func (s *session) get(table Table, opt getOptions, fields ...string) (string, bool, error) {
	if s.cache == nil {
		return s.getFromDatabase(table, fields...)
	}
	action, found, err := s.getFromCache(table, fields...)
	if err != nil {
		return action, found, err
	}
	if !found && opt.SyncFromDatabase {
		action, found, err = s.getFromDatabase(table, fields...)
		if err != nil {
			return action, found, err
		}
		if found {
			s.updateCache(table, fields...)
			// NOTE: this error sould be ignored
			return action_null, true, nil
		}
	}
	return action, found, err
}

func (s *session) getFromCache(table Table, fields ...string) (string, bool, error) {
	meta := table.Meta()
	tableKey := typeconv.ToString(table.Key())
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	fieldSize := len(fields)
	args := make([]string, 0, fieldSize)
	for _, field := range fields {
		args = append(args, s.cache.FieldName(s.eng.name, meta.Name(), tableKey, field))
	}
	values, err := s.cache.HMGet(s.cache.TableName(s.eng.name, meta.Name()), args...)
	if err != nil {
		return action_cache_hmget, false, err
	}
	if len(values) != fieldSize {
		return action_cache_hmget, false, ErrUnexpectedLength
	}
	found := false
	for i := 0; i < fieldSize; i++ {
		if values[i] != nil {
			if err := table.SetField(fields[i], typeconv.ToString(values[i])); err != nil {
				return action_set_field(meta.Name(), fields[i]), false, err
			}
			found = true
		}
	}
	return action_null, found, nil
}

func (s *session) getFromDatabase(table Table, fields ...string) (string, bool, error) {
	found, err := s.database.Get(table, fields...)
	if err != nil {
		return action_db_get, false, err
	}
	return action_null, found, nil
}

func (s *session) find(meta TableMeta, keys KeyList, setters FieldSetterList, fields []string) (string, error) {
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	_, action, err := s.findByFields(meta.Name(), keys, setters, FieldSlice(fields), nil)
	return action, err
}

func (s *session) clear(table string) (string, error) {
	if indexes, ok := s.eng.indexes[table]; ok {
		for _, index := range indexes {
			indexKey := s.cache.IndexKey(s.eng.name, index)
			if _, err := s.cache.Delete(indexKey); err != nil {
				return action_cache_del + ": delete index `" + indexKey + "`", err
			}
		}
	}
	key := s.cache.TableName(s.eng.name, table)
	if _, err := s.cache.Delete(key); err != nil {
		return action_cache_del, err
	}
	return action_null, nil
}

func (s *session) findByFields(table string, keys KeyList, setters FieldSetterList, fields FieldList, refs map[string]View) (map[string]StringKeys, string, error) {
	keySize := keys.Len()
	if keySize == 0 {
		return nil, action_null, nil
	}
	fieldSize := fields.Len()
	args := make([]string, 0, fieldSize*keySize)
	for i := 0; i < keySize; i++ {
		key := typeconv.ToString(keys.Key(i))
		for i := 0; i < fieldSize; i++ {
			args = append(args, s.cache.FieldName(s.eng.name, table, key, fields.Field(i)))
		}
	}
	values, err := s.cache.HMGet(s.cache.TableName(s.eng.name, table), args...)
	if err != nil {
		return nil, action_null, err
	}
	length := len(values)
	if length != fieldSize*keySize {
		return nil, action_null, ErrUnexpectedLength
	}
	var keysGroup map[string]StringKeys
	if len(refs) > 0 {
		keysGroup = make(map[string]StringKeys)
		for field := range refs {
			keysGroup[field] = StringKeys(make([]string, keySize))
		}
	}
	for i := 0; i+fieldSize <= length; i += fieldSize {
		index := i / fieldSize
		setter, err := setters.New(table, index, typeconv.ToString(keys.Key(index)))
		if err != nil {
			// NOTE: the error should be ignored
			continue
		}
		for j := 0; j < fieldSize; j++ {
			field := fields.Field(j)
			value := values[i+j]
			if value != nil {
				if err := setter.SetField(field, typeconv.ToString(value)); err != nil {
					return nil, action_set_field(table, field), err
				}
			}
			if ks, ok := keysGroup[field]; ok {
				if value == nil {
					ks[index] = ""
				} else {
					ks[index] = typeconv.ToString(value)
				}
				keysGroup[field] = ks
			}
		}
	}
	return keysGroup, action_null, nil
}

func (s *session) recursivelyLoadView(view View, keys KeyList, setters FieldSetterList) (string, error) {
	keysGroup, action, err := s.findByFields(view.Table(), keys, setters, view.Fields(), view.Refs())
	if err != nil {
		return action, err
	}
	refs := view.Refs()
	if refs == nil {
		return action_null, nil
	}
	if len(keysGroup) != len(refs) {
		return action, ErrUnexpectedLength
	}
	for field, ref := range refs {
		if tmpKeys, ok := keysGroup[field]; ok {
			if action, err := s.recursivelyLoadView(ref, tmpKeys, setters); err != nil {
				return action, err
			}
		} else {
			return action, ErrViewRefFieldMissing
		}
	}
	return action_null, nil
}

func (s *session) updateIndex(table ReadonlyTable, key interface{}, updatedFields []string) (action string, err error) {
	if indexes, ok := s.eng.indexes[table.Meta().Name()]; ok {
		for _, index := range indexes {
			if err = index.Update(s, table, key, updatedFields); err != nil {
				action = "update index `" + s.cache.IndexKey(s.eng.name, index) + "`"
				return
			}
		}
	}
	return
}

func (s *session) removeIndex(tableName string, keys ...interface{}) (action string, err error) {
	if indexes, ok := s.eng.indexes[tableName]; ok {
		for _, index := range indexes {
			if err = index.Remove(s, keys...); err != nil {
				action = "remove index `" + s.cache.IndexKey(s.eng.name, index) + "`"
				return
			}
		}
	}
	return
}

func (s *session) indexRank(index Index, key interface{}) (string, int64, error) {
	if indexRank, ok := index.(IndexRank); ok {
		rank, err := indexRank.Rank(key)
		return action_null, rank, err
	}
	rank, err := s.cache.ZRank(s.cache.IndexKey(s.eng.name, index), typeconv.ToString(key))
	if err != nil {
		if err == ErrNotFound {
			return action_null, InvalidRank, nil
		}
		return action_cache_zrank, InvalidRank, err
	}
	if rank < 0 {
		rank = InvalidRank
	}
	return action_null, rank, nil
}

func (s *session) indexScore(index Index, key interface{}) (string, int64, error) {
	if indexScore, ok := index.(IndexScore); ok {
		score, err := indexScore.Score(key)
		return action_null, score, err
	}
	score, err := s.cache.ZScore(s.cache.IndexKey(s.eng.name, index), typeconv.ToString(key))
	if err != nil {
		if err == ErrNotFound {
			return action_null, InvalidScore, nil
		}
		return action_cache_zscore, score, err
	}
	return action_null, score, nil
}
