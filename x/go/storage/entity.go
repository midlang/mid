package storage

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

// FieldGetter get value by field
type FieldGetter interface {
	GetField(field string) (interface{}, bool)
}

// FieldGetter set value by field
type FieldSetter interface {
	SetField(field, value string) error
}

// TableMeta holds table meta information
type TableMeta interface {
	// Name returns name of table
	Name() string
	// Key returns name of key field
	Key() string
	// Fields returns names of all fields except key field
	Fields() []string
}

//-------------------
// Compose interface
//-------------------

// ReadonlyTable represents a read-only table
type ReadonlyTable interface {
	Meta() TableMeta
	Key() interface{}
	FieldGetter
}

// Table represents a table in sql database, and hash table in nosql database
type Table interface {
	ReadonlyTable
	SetKey(string) error
	FieldSetter
}

// TableListContainer is a linear container holds and creates Table
type TableListContainer interface {
	Len() int
	New(tableName string, index int, key string) (Table, error)
}

// View represents a view references a table
type View interface {
	Table() string
	Fields() FieldList
	Refs() map[string]View
}

// Index represents a sorted set for field of table
type Index interface {
	Name() string
	Table() string
	Update(s Session, table ReadonlyTable, key interface{}, updatedFields []string) error
	Remove(s Session, keys ...interface{}) error
}
