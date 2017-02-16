package storage

import (
	"errors"
	"strconv"
	"testing"

	"github.com/mkideal/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	redisc *mockRedisclient

	errArgumentsLength = errors.New("invalid arguments length")
	errUnknownField    = errors.New("unknown field")
)

func init() {
	redisc = &mockRedisclient{
		tables: map[string]mockTable{},
	}
}

type mockTable map[string]interface{}

func (t mockTable) set(field string, value interface{}) {
	t[field] = value
}

func (t mockTable) get(field string) (string, bool) {
	value, found := t[field]
	return ToString(value), found
}

func (t mockTable) del(field string) bool {
	if _, found := t[field]; found {
		delete(t, field)
		return true
	}
	return false
}

type mockRedisclient struct {
	tables map[string]mockTable
}

func (c *mockRedisclient) getTable(name string) mockTable {
	table, ok := c.tables[name]
	if ok {
		return table
	}
	table = mockTable{}
	c.tables[name] = table
	return table
}

func (c *mockRedisclient) HsetMulti(args ...interface{}) (string, error) {
	if len(args)%2 == 0 || len(args) == 1 {
		return "", errArgumentsLength
	}
	tableName := ToString(args[0])
	table := c.getTable(tableName)
	for i := 1; i+1 < len(args); i += 2 {
		table.set(ToString(args[i]), args[i+1])
	}
	return ToString(len(args) / 2), nil
}

func (c *mockRedisclient) Hmgetstrings(args ...interface{}) (int, []*string, error) {
	if len(args) < 2 {
		return 0, nil, errArgumentsLength
	}
	values := make([]*string, 0, len(args)-1)
	tableName := ToString(args[0])
	table := c.getTable(tableName)
	n := 0
	for i := 1; i < len(args); i++ {
		field := ToString(args[i])
		value, found := table.get(field)
		if found {
			values = append(values, &value)
		} else {
			values = append(values, nil)
			n++
		}
	}
	return n, values, nil
}

func (c *mockRedisclient) HdelMulti(args ...interface{}) (string, error) {
	if len(args) < 2 {
		return "0", errArgumentsLength
	}
	tableName := ToString(args[0])
	table := c.getTable(tableName)
	n := 0
	for i := 1; i < len(args); i++ {
		field := ToString(args[i])
		if table.del(field) {
			n++
		}
	}
	return ToString(n), nil
}

func (c *mockRedisclient) Delete(key string) error {
	delete(c.tables, key)
	return nil
}

func (c *mockRedisclient) ZaddMultiScore(key string, member interface{}, score int64) (bool, error) {
	return false, nil
}

func (c *mockRedisclient) ZremMulti(key_members ...interface{}) error {
	return nil
}

func (c *mockRedisclient) Zrank(key string, member interface{}) int { return 0 }
func (c *mockRedisclient) Zscore64(key string, member interface{}) (int64, bool, error) {
	return 0, false, nil
}

func TestOrm(t *testing.T) {
	defer log.Uninit(log.InitConsole(log.LvFATAL))

	eng := NewEngine("test", redisc)
	eng.SetErrorHandler(func(action string, err error) error {
		log.Printf(ErrorHandlerDepth, log.LvWARN, "<%s>: %v", action, err)
		return err
	})

	// Insert
	inserted := &User{Id: 1, Name: "test1", Age: 10}
	eng.Insert(inserted)
	t.Logf("insert 1: %v", redisc.tables)

	// Get
	loaded := &User{Id: 1}
	eng.Get(loaded)
	assert.Equal(t, inserted.Id, loaded.Id)
	assert.Equal(t, inserted.Name, loaded.Name)
	assert.Equal(t, inserted.Age, loaded.Age)

	// Update
	inserted.Age = 20
	eng.Update(inserted, "age")
	loaded = &User{Id: 1}
	eng.Get(loaded)
	assert.Equal(t, inserted.Id, loaded.Id)
	assert.Equal(t, inserted.Name, loaded.Name)
	assert.Equal(t, inserted.Age, loaded.Age)
	assert.Equal(t, 20, loaded.Age)

	// Remove
	eng.Remove(inserted)
	found, err := eng.Get(&User{Id: 1})
	assert.Nil(t, err)
	assert.False(t, found)

	users := []*User{
		{Id: 1, Name: "test1", Age: 10},
		{Id: 2, Name: "test2", Age: 20},
		{Id: 3, Name: "test3", Age: 30},
	}
	keys := make([]int64, 0, len(users))
	for _, u := range users {
		eng.Insert(u)
		keys = append(keys, u.Id)
	}
	t.Logf("insert %v: %v", keys, redisc.tables)

	// Find
	tf := newUserSlice(len(users))
	eng.Find(userMeta, Int64Keys(keys), tf)
	require.Equal(t, len(users), len(tf.data))
	for i := 0; i < len(users); i++ {
		assert.Equal(t, users[i].Id, tf.data[i].Id)
		assert.Equal(t, users[i].Name, tf.data[i].Name)
		assert.Equal(t, users[i].Age, tf.data[i].Age)
	}

	// Error
	loaded = &User{Id: 1}
	err = eng.Update(loaded, "invalid_field")
	require.NotNil(t, err)
	t.Logf(err.Error())
}

type UserMeta struct{}

var userMeta = UserMeta{}

func (meta UserMeta) Name() string     { return "user" }
func (meta UserMeta) Fields() []string { return []string{"name", "age"} }

type userSlice struct {
	data []User
}

func newUserSlice(cap int) *userSlice {
	return &userSlice{
		data: make([]User, 0, cap),
	}
}

func (us *userSlice) New(table string, index int, key string) (FieldSetter, error) {
	for len(us.data) <= index {
		us.data = append(us.data, User{})
	}
	err := us.data[index].SetKey(key)
	return &us.data[index], err
}

type User struct {
	Id   int64
	Name string
	Age  int
}

func (u User) Meta() TableMeta  { return userMeta }
func (u User) Key() interface{} { return u.Id }
func (u *User) SetKey(value string) error {
	return setInt64(&u.Id, value)
}

func (u User) GetField(field string) (interface{}, bool) {
	switch field {
	case "name":
		return u.Name, true
	case "age":
		return u.Age, true
	default:
		return nil, false
	}
}

func (u *User) SetField(field string, value string) error {
	switch field {
	case "name":
		u.Name = value
	case "age":
		return setInt(&u.Age, value)
	default:
		return errUnknownField
	}
	return nil
}

func setInt(ptr *int, value string) error {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*ptr = int(i)
	return nil
}

func setInt64(ptr *int64, value string) error {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	*ptr = i
	return nil
}

type ViewUserName struct {
	Name string
}

func (ViewUserName) Table() string         { return userMeta.Name() }
func (ViewUserName) Fields() FieldList     { return Field("name") }
func (ViewUserName) Refs() map[string]View { return nil }

func TestOrmView(t *testing.T) {
	defer log.Uninit(log.InitConsole(log.LvFATAL))

	eng := NewEngine("test", redisc)
	eng.SetErrorHandler(func(action string, err error) error {
		log.Printf(ErrorHandlerDepth, log.LvWARN, "<%s>: %v", action, err)
		return err
	})

	// Insert
	inserted := &User{Id: 1, Name: "test1", Age: 10}
	eng.Insert(inserted)
	t.Logf("insert 1: %v", redisc.tables)

	// LoadView
	view := ViewUserName{}
	fs := newUserSlice(1)
	keys := Int64Keys([]int64{inserted.Id})
	eng.FindView(view, keys, fs)
	t.Logf("insert 1: %v", fs)
}
