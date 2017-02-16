package storage

const (
	action_null      = ""
	action_hmget     = "redis.HMGET"
	action_hdel      = "redis.HDEL"
	action_del       = "redis.DEL"
	action_zscore    = "redis.ZSCORE"
	action_zrange    = "redis.ZRANGE"
	action_zrevrange = "redis.ZREVRANGE"
)

func action_get_field(table, field string) string {
	return "table `" + table + "` GetField `" + field + "`"
}
func action_set_field(table, field string) string {
	return "table `" + table + "` SetField `" + field + "`"
}

// Engine is the core engine of redis orm package
type Engine struct {
	name    string
	redisc  RedisClient
	indexes map[string]map[string]Index

	errorHandler ErrorHandler
}

func NewEngine(name string, redisc RedisClient) *Engine {
	return &Engine{
		name:    name,
		redisc:  redisc,
		indexes: make(map[string]map[string]Index),
	}
}

func (eng *Engine) catch(action string, err error) error {
	if err != nil && eng.errorHandler != nil {
		err = eng.errorHandler(action, err)
	}
	return err
}

func (eng *Engine) tableName(name string) string {
	return eng.name + "@" + name
}

func (eng *Engine) fieldName(key interface{}, field string) string {
	return ToString(key) + ":" + field
}

//-------------
// engine APIs
//-------------

// Name returns database name
func (eng *Engine) Name() string { return eng.name }

// RedisClient returns redis client
func (eng *Engine) RedisClient() RedisClient { return eng.redisc }

// IndexKey calculates redis key of an index
func (eng *Engine) IndexKey(index Index) string {
	return eng.name + "@" + index.Table() + ":" + index.Name()
}

// SetErrorHandler sets handler for handling error
func (eng *Engine) SetErrorHandler(eh ErrorHandler) {
	eng.errorHandler = eh
}

// CreateIndex creates an index
func (eng *Engine) CreateIndex(index Index) {
	tableName := index.Table()
	idx, ok := eng.indexes[tableName]
	if !ok {
		idx = make(map[string]Index)
		eng.indexes[tableName] = idx
	}
	idx[index.Name()] = index
}

// Insert inserts new records or updates all fields of records
func (eng *Engine) Insert(tables ...ReadonlyTable) error {
	for _, table := range tables {
		action, err := eng.update(table)
		if err != nil {
			return eng.catch("Insert: "+action, err)
		}
	}
	return nil
}

// Update updates specific fields of record
func (eng *Engine) Update(table ReadonlyTable, fields ...string) error {
	action, err := eng.update(table, fields...)
	if err != nil {
		return eng.catch("Update: "+action, err)
	}
	return nil
}

// BatchUpdate updates specific fields of records
func (eng *Engine) BatchUpdate(tables ReadonlyTableList, fields ...string) error {
	length := tables.Len()
	for i := 0; i < length; i++ {
		table := tables.ReadonlyTable(i)
		action, err := eng.update(table, fields...)
		if err != nil {
			return eng.catch("BatchUpdate: "+action, err)
		}
	}
	return nil
}

// Find gets many records
func (eng *Engine) Find(meta TableMeta, keys KeyList, setters FieldSetterList, fields ...string) error {
	action, err := eng.find(meta, keys, setters, fields)
	if err != nil {
		return eng.catch("Find: "+action, err)
	}
	return nil
}

// Get gets one record by specific fields. It will gets all fields if fields is empty
func (eng *Engine) Get(table WriteonlyTable, fields ...string) (bool, error) {
	action, ok, err := eng.get(table, fields...)
	if err != nil {
		return ok, eng.catch("Get: "+action, err)
	}
	return ok, nil
}

// Remove removes one record
func (eng *Engine) Remove(table ReadonlyTable) error {
	meta := table.Meta()
	action, err := eng.remove(meta, table.Key())
	if err != nil {
		return eng.catch("Remove: "+action, err)
	}
	return nil
}

// RemoveKeys removes records by keys
func (eng *Engine) RemoveKeys(meta TableMeta, keys ...interface{}) error {
	action, err := eng.remove(meta, keys...)
	if err != nil {
		return eng.catch("RemoveKeys: "+action, err)
	}
	return nil
}

// DropTable removes whole table
func (eng *Engine) DropTable(table string) error {
	key := eng.tableName(table)
	if indexes, ok := eng.indexes[table]; ok {
		for _, index := range indexes {
			indexKey := eng.IndexKey(index)
			if err := eng.redisc.Delete(indexKey); err != nil {
				return eng.catch("DropTable "+table+": "+action_del+": delete index `"+indexKey+"`", err)
			}
		}
	}
	if err := eng.redisc.Delete(key); err != nil {
		eng.catch("DropTable "+table+": "+action_del, err)
	}
	return nil
}

// FindView loads view by keys and store loaded data to setters
func (eng *Engine) FindView(view View, keys KeyList, setters FieldSetterList) error {
	action, err := eng.recursivelyLoadView(view, keys, setters)
	if err != nil {
		return eng.catch("FindView: "+action, err)
	}
	return nil
}

// IndexRank gets rank of table key in index, returns InvalidRank if key not found
func (eng *Engine) IndexRank(index Index, key interface{}) (int, error) {
	if indexRank, ok := index.(IndexRank); ok {
		rank, err := indexRank.Rank(key)
		if err != nil {
			return rank, eng.catch("IndexRank: ", err)
		}
		return rank, nil
	}
	action, rank, err := eng.indexRank(index, key)
	if err != nil {
		return rank, eng.catch("IndexRank: "+action, err)
	}
	return rank, nil
}

// IndexScore gets score of table key in index, returns InvalidScore if key not found
func (eng *Engine) IndexScore(index Index, key interface{}) (int64, error) {
	if indexScore, ok := index.(IndexScore); ok {
		score, err := indexScore.Score(key)
		if err != nil {
			return score, eng.catch("IndexScore: ", err)
		}
		return score, nil
	}
	action, score, err := eng.indexScore(index, key)
	if err != nil {
		return score, eng.catch("IndexScore: "+action, err)
	}
	return score, nil
}

//----------------
// implementation
//----------------

func (eng *Engine) update(table ReadonlyTable, fields ...string) (string, error) {
	var (
		meta = table.Meta()
		key  = table.Key()
	)
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	args := make([]interface{}, 0, len(fields)*2+1)
	args = append(args, eng.tableName(meta.Name()))
	for _, field := range fields {
		args = append(args, eng.fieldName(key, field))
		value, ok := table.GetField(field)
		if !ok {
			return action_get_field(meta.Name(), field), ErrFieldNotFound
		}
		args = append(args, value)
	}
	action, err := eng.updateIndex(table, key, fields)
	if err != nil {
		return action, err
	}
	_, err = eng.redisc.HsetMulti(args...)
	return action_hmget, err
}

func (eng *Engine) remove(meta TableMeta, keys ...interface{}) (string, error) {
	fields := meta.Fields()
	args := make([]interface{}, 0, len(fields)+1)
	args = append(args, eng.tableName(meta.Name()))
	for _, key := range keys {
		for _, field := range fields {
			args = append(args, eng.fieldName(key, field))
		}
	}
	if action, err := eng.removeIndex(meta.Name(), keys...); err != nil {
		return action, err
	}
	_, err := eng.redisc.HdelMulti(args...)
	return action_hdel, err
}

func (eng *Engine) get(table WriteonlyTable, fields ...string) (string, bool, error) {
	meta := table.Meta()
	tableKey := ToString(table.Key())
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	fieldSize := len(fields)
	args := make([]interface{}, 0, fieldSize+1)
	args = append(args, eng.tableName(meta.Name()))
	for _, field := range fields {
		args = append(args, eng.fieldName(tableKey, field))
	}
	_, values, err := eng.redisc.Hmgetstrings(args...)
	if err != nil {
		return action_hmget, false, err
	}
	if len(values) != fieldSize {
		return action_hmget, false, ErrUnexpectedLength
	}
	found := false
	for i := 0; i < fieldSize; i++ {
		if values[i] != nil {
			if err := table.SetField(fields[i], *values[i]); err != nil {
				return action_set_field(meta.Name(), fields[i]), false, err
			}
			found = true
		}
	}
	return action_null, found, nil
}

func (eng *Engine) find(meta TableMeta, keys KeyList, setters FieldSetterList, fields []string) (string, error) {
	if len(fields) == 0 {
		fields = meta.Fields()
	}
	_, action, err := eng.findByFields(meta.Name(), keys, setters, FieldSlice(fields), nil)
	return action, err
}

func (eng *Engine) findByFields(table string, keys KeyList, setters FieldSetterList, fields FieldList, refs map[string]View) (map[string]StringKeys, string, error) {
	keySize := keys.Len()
	if keySize == 0 {
		return nil, action_null, nil
	}
	fieldSize := fields.Len()
	args := make([]interface{}, 0, fieldSize*keySize+1)
	args = append(args, eng.tableName(table))
	for i := 0; i < keySize; i++ {
		key := ToString(keys.Key(i))
		for i := 0; i < fieldSize; i++ {
			args = append(args, eng.fieldName(key, fields.Field(i)))
		}
	}
	_, values, err := eng.redisc.Hmgetstrings(args...)
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
		setter, err := setters.New(table, index, ToString(keys.Key(index)))
		if err != nil {
			// NOTE: the error should be ignored
			continue
		}
		for j := 0; j < fieldSize; j++ {
			field := fields.Field(j)
			value := values[i+j]
			if value != nil {
				if err := setter.SetField(field, *value); err != nil {
					return nil, action_set_field(table, field), err
				}
			}
			if ks, ok := keysGroup[field]; ok {
				if value == nil {
					ks[index] = ""
				} else {
					ks[index] = *value
				}
				keysGroup[field] = ks
			}
		}
	}
	return keysGroup, action_null, nil
}

func (eng *Engine) recursivelyLoadView(view View, keys KeyList, setters FieldSetterList) (string, error) {
	keysGroup, action, err := eng.findByFields(view.Table(), keys, setters, view.Fields(), view.Refs())
	if err != nil {
		return action, err
	}
	refs := view.Refs()
	if refs == nil {
		return "", nil
	}
	if len(keysGroup) != len(refs) {
		return action, ErrUnexpectedLength
	}
	for field, ref := range refs {
		if tmpKeys, ok := keysGroup[field]; ok {
			if action, err := eng.recursivelyLoadView(ref, tmpKeys, setters); err != nil {
				return action, err
			}
		} else {
			return action, ErrViewRefFieldMissing
		}
	}
	return "", nil
}

func (eng *Engine) updateIndex(table ReadonlyTable, key interface{}, updatedFields []string) (action string, err error) {
	if indexes, ok := eng.indexes[table.Meta().Name()]; ok {
		for _, index := range indexes {
			if err = index.Update(eng, table, key, updatedFields); err != nil {
				action = "update index `" + eng.IndexKey(index) + "`"
				return
			}
		}
	}
	return
}

func (eng *Engine) removeIndex(tableName string, keys ...interface{}) (action string, err error) {
	if indexes, ok := eng.indexes[tableName]; ok {
		for _, index := range indexes {
			if err = index.Remove(eng, keys...); err != nil {
				action = "remove index `" + eng.IndexKey(index) + "`"
				return
			}
		}
	}
	return
}

func (eng *Engine) indexRank(index Index, key interface{}) (string, int, error) {
	rank := eng.redisc.Zrank(eng.IndexKey(index), key)
	if rank < 0 {
		rank = InvalidRank
	}
	return "", rank, nil
}

func (eng *Engine) indexScore(index Index, key interface{}) (string, int64, error) {
	score, found, err := eng.redisc.Zscore64(eng.IndexKey(index), key)
	if err != nil {
		return action_zscore, score, err
	}
	if !found {
		return "", InvalidScore, nil
	}
	return "", score, nil
}
