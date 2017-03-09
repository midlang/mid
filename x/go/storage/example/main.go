package main

import (
	"flag"
	"time"

	"github.com/mkideal/log"
	"gopkg.in/redis.v5"

	"github.com/midlang/mid/x/go/storage"
	"github.com/midlang/mid/x/go/storage/example/demo"
	"github.com/midlang/mid/x/go/storage/goredisproxy"
)

var (
	flRedisHost = flag.String("redis_host", "127.0.0.1:6379", "redis host")
	flRedisPwd  = flag.String("redis_pwd", "", "redis password")
)

func main() {
	flag.Parse()
	defer log.Uninit(log.InitColoredConsole(log.LvFATAL))
	client := redis.NewClient(&redis.Options{
		Addr:         *flRedisHost,
		Password:     *flRedisPwd,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	client.FlushDb()

	cache := goredisproxy.New(client)

	// `数据库名`的含义是指在redis的每个key前面会自动添加 `<数据库名>@` 前缀
	// 比如指定为 test 后，一个名为 abc 的表的 key 就是 test@abc
	eng := storage.NewEngine("test", storage.NullDatabaseProxy, cache) // 这里不存数据库,使用了 NullDatabaseProxy
	eng.AddIndex(demo.UserAgeIndexVar)
	eng.AddIndex(demo.UserAddrIndexVar)
	eng.AddIndex(demo.ProductPriceIndexVar)

	// SetErrorHandler 设置一个错误回调函数,当 eng 执行遇到错误时就会调用此回调函数
	// 回调函数的类型为 `func(action string, err error) error`
	// `action` 参数指明在进行什么操作时出的错
	// `err` 参数为错误值
	// 返回值也是一个error,用于替换传入的错误
	// 下面这个回调函数只是打印了一条日志,然后将错误原样返回
	// storage.ErrorHandlerDepth 常量值指明这个回调函数被调用的栈深度,方便日志打印时定位到接口调用的地方
	eng.SetErrorHandler(func(action string, err error) error {
		log.Printf(storage.ErrorHandlerDepth, log.LvWARN, "<%s>: %v", action, err)
		return err
	})

	api := eng.NewSession()
	defer api.Close() // 也可以不用NewSession，直接使用 eng,试一试 api := eng

	// Insert 插入任意多条记录到给定的表中,Insert 函数的入参必须实现 storage.ReadonlyTable 接口
	// 此函数接收任意多个参数
	inserted := &demo.User{Id: 1, Name: "test1", Age: 10}
	api.Insert(inserted)

	// Get 获取一条记录,Get 函数的第一个参数必须实现 storage.Table 接口
	// Get 方法只能根据 WriteonlyTable 的 Key 来获取值,不能根据其他字段
	// 后面接收字符串类型的可变参数,这些参数指明了需要获取的 WriteonlyTable 的字段名
	// 如果没有参数,则获取 WriteonlyTable 的所有字段值
	// 比如下面的
	//
	//		api.Get(loaded)
	//
	// 根据User的key(这里就是user的Id) 获取User的全部字段
	// 如果写成
	//
	//		api.Get(loaded, demo.UserMetaVar.F_age, demo.UserMetaVar.F_name)
	//
	// 则获取User的Age和Name
	loaded := &demo.User{Id: 1}
	found, err := api.Get(loaded)
	log.Info("Get: found=%v, error=%v, value=%v", found, err, loaded)

	// Update 根据 Table 的 Key 更新全部或指定字段值
	// Update 方法的第一个参数必须实现 storage.ReadonlyTable 接口
	// 后面接收字符串类型的可变参数,这些参数指明了需要更新的 ReadonlyTable 的字段名
	// 如果没有参数,则更新 ReadonlyTable 的所有字段值
	// 比如下面的
	//
	//		api.Update(inserted, demo.UserMetaVar.F_age)
	//
	// 根据User的Id更新Age字段
	// 如果写成
	//
	//		api.Update(inserted)
	//
	// 则更新全部字段
	inserted.Age = 20
	api.Update(inserted, demo.UserMetaVar.F_age)
	//loaded = &demo.User{Id: 1}
	//found, err = api.Get(loaded)
	//log.Info("Get: found=%v, error=%v, value=%v", found, err, loaded)

	// Remove 删除一条记录
	// 第一个参数必须实现 storage.ReadonlyTable 接口
	// 除了 Remove,还有 RemoveKeys 和 DropTable 方法可以用于删除数据
	// RemoveKeys 根据一组 key 删除一个表的一组记录, DropTable 删除整个表
	api.Remove(inserted) // 等同于 api.RemoveRecords(inserted.Meta(), inserted.Key())
	//found, err = api.Get(&demo.User{Id: inserted.Id})
	//log.Info("Remove: found=%v, error=%v", found, err)

	users := []demo.User{
		{Id: 1, Name: "test1", Age: 10, AddrId: 1000, ProductId: 100},
		{Id: 2, Name: "test2", Age: 20, AddrId: 2000, ProductId: 200},
		{Id: 3, Name: "test3", Age: 30, AddrId: 3000, ProductId: 300},
	}
	products := []demo.Product{
		{Id: 100, Price: 1, Name: "p1", Image: "img1", Desc: "desc1"},
		{Id: 200, Price: 2, Name: "p2", Image: "img2", Desc: "desc2"},
		{Id: 300, Price: 3, Name: "p3", Image: "img3", Desc: "desc3"},
	}
	addresses := []demo.Address{
		{Id: 1000, Addr: "Beijing"},
		{Id: 2000, Addr: "Shanghai"},
		{Id: 3000, Addr: "Taiwan"},
	}
	keys := make([]int64, 0, len(users))
	for i := range users {
		api.Insert(&users[i])
		keys = append(keys, users[i].Id)
	}
	for i := range products {
		api.Insert(&products[i])
	}
	for i := range addresses {
		api.Insert(&addresses[i])
	}

	// Find 根据一组 key 查找一组记录
	// 第一个参数必须实现 storage.TableMeta 接口
	// 第二个参数必须实现 storage.KeyList 接口(storage内置IntKeys,Int64Keys,Uint64Keys,StringKeys,InterfaceKeys等封装数组为KeyList的类型)
	// 第三个参数必须实现 storage.FieldSetterList 接口
	// 再后面接收字符串类型的可变参数,用于指定需要获取的字段名,不传参数则获取表的所有字段
	us := demo.NewUserSlice(len(users))
	api.Find(demo.UserMetaVar, storage.Int64Keys(keys), us)
	log.Info("Find users: %v", us.Slice())

	// 测试错误处理回调,这里尝试 Update 一个无效字段名
	// 如果redis配置正确,运行 `go run main.go` 可以看到如下的错误日志:
	//
	// [W 2016/12/05 16:37:18.708 main.go:154] <Update: table `user` GetField `invalid_field`>: field not found
	//
	// 这条错误信息表明在 Update 函数中,user 表在获取字段 `invalid_field` 时发现这个字段找不到
	// 注意,这并不是说 redis 中没有这个字段,而是说 `user` 表的定义中没有这个字段
	// 日志被定位到 main.go:157 即 Update 函数被调用的地方
	loaded = &demo.User{Id: 1}
	err = api.Update(loaded, "invalid_field")

	// FindView 按视图获取数据
	// 在 User 中,AddrId 字段引用了 Address 表,ProductId 字段引用了 Product 表
	// 而 Product 中的 SkuId 字段又引用了 Sku 表
	// 所以现在 UserView 视图就是 User 和 Address,Product,Sku 的联合视图
	// 第一个参数必须实现 storage.View 接口
	// 第二个参数必须实现 storage.KeyList 接口
	// 第三个参数必须实现 storage.FieldSetterList 接口
	userview := demo.NewUserViewSlice(len(users))
	api.FindView(demo.UserViewVar, storage.Int64Keys(keys), userview)
	log.WithJSON(userview.Slice()).Info("FindView: userview")

	// IndexRank 获取索引的排名
	rank, _ := api.IndexRank(demo.UserAgeIndexVar, 1)
	log.Info("IndexRank: user %d UserAgeIndex: %d", 1, rank)
	rank, _ = api.IndexRank(demo.UserAgeIndexVar, -1)
	log.Info("IndexRank: user %d UserAgeIndex equal to storage.InvalidRank? %v", -1, rank == storage.InvalidRank)

	// IndexScore 获取索引的score
	score, _ := api.IndexScore(demo.UserAgeIndexVar, 1)
	log.Info("IndexScore: user %d UserAgeIndex: %d", 1, score)
	score, _ = api.IndexScore(demo.UserAgeIndexVar, -1)
	log.Info("IndexScore: user %d UserAgeIndex equal to storage.InvalidScore? %v", -1, score == storage.InvalidScore)

	// Clear 删除表,对应的索引也会被清除
	api.Clear(demo.UserMetaVar.Name())
	api.Clear(demo.AddressMetaVar.Name())
	api.Clear(demo.ProductMetaVar.Name())
}
