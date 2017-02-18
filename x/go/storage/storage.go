package storage

import (
	"math"
)

const (
	Unused = 0

	ErrorHandlerDepth = 4

	InvalidRank  = -1
	InvalidScore = math.MinInt64
)

type RedisClient interface {
	Delete(key string) error
	HsetMulti(...interface{}) (string, error)
	Hmgetstrings(...interface{}) (int, []*string, error)
	HdelMulti(...interface{}) (string, error)
	ZaddMultiScore(key string, member interface{}, score int64) (bool, error)
	ZremMulti(key_members ...interface{}) error
	Zrank(key string, member interface{}) int
	Zscore64(key string, member interface{}) (int64, bool, error)
}

// KeyList holds n keys
type KeyList interface {
	Len() int
	Key(int) interface{}
}

type IntKeys []int
type Int8Keys []int8
type Int16Keys []int16
type Int32Keys []int32
type Int64Keys []int64
type UintKeys []uint
type Uint8Keys []uint8
type Uint16Keys []uint16
type Uint32Keys []uint32
type Uint64Keys []uint64
type StringKeys []string
type InterfaceKeys []interface{}

func (keys IntKeys) Len() int       { return len(keys) }
func (keys Int8Keys) Len() int      { return len(keys) }
func (keys Int16Keys) Len() int     { return len(keys) }
func (keys Int32Keys) Len() int     { return len(keys) }
func (keys Int64Keys) Len() int     { return len(keys) }
func (keys UintKeys) Len() int      { return len(keys) }
func (keys Uint8Keys) Len() int     { return len(keys) }
func (keys Uint16Keys) Len() int    { return len(keys) }
func (keys Uint32Keys) Len() int    { return len(keys) }
func (keys Uint64Keys) Len() int    { return len(keys) }
func (keys StringKeys) Len() int    { return len(keys) }
func (keys InterfaceKeys) Len() int { return len(keys) }

func (keys IntKeys) Key(i int) interface{}       { return keys[i] }
func (keys Int8Keys) Key(i int) interface{}      { return keys[i] }
func (keys Int16Keys) Key(i int) interface{}     { return keys[i] }
func (keys Int32Keys) Key(i int) interface{}     { return keys[i] }
func (keys Int64Keys) Key(i int) interface{}     { return keys[i] }
func (keys UintKeys) Key(i int) interface{}      { return keys[i] }
func (keys Uint8Keys) Key(i int) interface{}     { return keys[i] }
func (keys Uint16Keys) Key(i int) interface{}    { return keys[i] }
func (keys Uint32Keys) Key(i int) interface{}    { return keys[i] }
func (keys Uint64Keys) Key(i int) interface{}    { return keys[i] }
func (keys StringKeys) Key(i int) interface{}    { return keys[i] }
func (keys InterfaceKeys) Key(i int) interface{} { return keys[i] }

// FieldList holds n fields
type FieldList interface {
	Len() int
	Field(int) string
}

// Field implements FieldList which atmost contains one value
type Field string

func (f Field) Len() int {
	if f == "" {
		return 0
	}
	return 1
}

func (f Field) Field(i int) string { return string(f) }

// FieldSlice implements FieldList
type FieldSlice []string

func (fs FieldSlice) Len() int           { return len(fs) }
func (fs FieldSlice) Field(i int) string { return fs[i] }

//-----------------
// Basic interface
//-----------------

type FieldGetter interface {
	GetField(field string) (interface{}, bool)
}

type FieldSetter interface {
	SetField(field, value string) error
}

// TableMeta holds table meta information
type TableMeta interface {
	Name() string
	Key() string
	Fields() []string
}

//-------------------
// Compose interface
//-------------------

type TableInfo interface {
	Meta() TableMeta
	Key() interface{}
}

type ReadonlyTable interface {
	TableInfo
	FieldGetter
}

type WriteonlyTable interface {
	TableInfo
	FieldSetter
	SetKey(string) error
}

type Table interface {
	TableInfo
	FieldGetter
	FieldSetter
	SetKey(string) error
}

type FieldSetterList interface {
	New(table string, index int, key string) (FieldSetter, error)
}

type ReadonlyTableList interface {
	Len() int
	ReadonlyTable(i int) ReadonlyTable
}

type View interface {
	Table() string
	Fields() FieldList
	Refs() map[string]View
}

type Index interface {
	Name() string
	Table() string
	Update(eng *Engine, table ReadonlyTable, key interface{}, updatedFields []string) error
	Remove(eng *Engine, keys ...interface{}) error
}

type IndexRank interface {
	Rank(key interface{}) (int, error)
}

type IndexScore interface {
	Score(key interface{}) (int64, error)
}

type IndexRanger interface {
	Range(eng *Engine, start, end int) (KeyList, error)
}

type IndexScoreRanger interface {
	RangeByScore(eng *Engine, min, max int64) (KeyList, error)
}

type IndexRevRanger interface {
	RevRange(eng *Engine, start, end int) (KeyList, error)
}

type IndexScoreRevRanger interface {
	RevRangeByScore(eng *Engine, max, min int64) (KeyList, error)
}

type ErrorHandler func(action string, err error) error
