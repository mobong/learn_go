package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"learn_go/conf"
	"strconv"
	"time"
)

func main() {
	client, ctx := startRedis()
	strOperate(ctx, client)
	hashOperate(ctx, client)
	listOperate(ctx, client)
	setOperate(ctx, client)
	sortedSetOperate(ctx, client)
	transactionOperate(ctx, client)
	closeRedis(client)
}

// 连接redis
func startRedis() (*redis.Client, context.Context) {
	conf.ConfInit()
	host := viper.GetString("redis.host")
	port := viper.GetInt("redis.port")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	addr := host + ":" + strconv.Itoa(port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return client, ctx
}

// 关闭redis
func closeRedis(client *redis.Client) {
	err := client.Close()
	if err != nil {
		fmt.Println(err)
	}
}

// 字符串操作
func strOperate(ctx context.Context, client *redis.Client) {
	// 设置字符串
	err := client.Set(ctx, "mykey", "myvalue", 0).Err()
	if err != nil {
		fmt.Println(err)
	}
	// 获取字符串
	strValue, err := client.Get(ctx, "mykey").Result()
	if err == redis.Nil {
		fmt.Println("key不存在")
	} else if err != nil {
		fmt.Println("strValue error:", err)
	} else {
		fmt.Println("mykey:", strValue)
	}
	// 删除字符串
	err = client.Del(ctx, "mykey").Err()
	if err != nil {
		fmt.Println(err)
	}

	// 设置带有过期时间的键值对
	err = client.Set(ctx, "mykey2", "myvalue2", 10*time.Second).Err()
	if err != nil {
		fmt.Println(err)
	}

	// 获取键的剩余过期时间
	ttl, err := client.TTL(ctx, "mykey2").Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("mykey2剩余过期时间:", ttl)
	}
}

// 哈希操作
func hashOperate(ctx context.Context, client *redis.Client) {
	// 设置哈希字段
	err := client.HSet(ctx, "myhash", "field", "value").Err()
	if err != nil {
		fmt.Println(err)
	}

	// 获取哈希字段的值
	val, err := client.HGet(ctx, "myhash", "field").Result()
	if err == redis.Nil {
		fmt.Println("字段不存在")
	} else if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("field:", val)
	}

	// 删除哈希字段
	err = client.HDel(ctx, "myhash", "field").Err()
	if err != nil {
		fmt.Println(err)
	}
}

// 列表操作
func listOperate(ctx context.Context, client *redis.Client) {
	// 在列表尾部插入元素
	err := client.RPush(ctx, "mylist", "element1", "element2").Err()
	if err != nil {
		fmt.Println(err)
	}

	// 获取列表指定范围的元素
	elements, err := client.LRange(ctx, "mylist", 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("列表元素:", elements)
	}

	// 删除列表中的元素
	err = client.LRem(ctx, "mylist", 0, "element1").Err()
	if err != nil {
		fmt.Println(err)
	}
}

// 集合操作
func setOperate(ctx context.Context, client *redis.Client) {
	// 添加集合元素
	err := client.SAdd(ctx, "myset", "element1", "element2").Err()
	if err != nil {
		fmt.Println(err)
	}

	// 获取集合所有元素
	elements, err := client.SMembers(ctx, "myset").Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("集合元素:", elements)
	}

	// 删除集合元素
	err = client.SRem(ctx, "myset", "element1", "element2").Err()
	if err != nil {
		fmt.Println(err)
	}
}

// 有序集合操作
func sortedSetOperate(ctx context.Context, client *redis.Client) {
	// 添加有序集合元素
	err := client.ZAdd(ctx, "mysorted-set", &redis.Z{Score: 1, Member: "element1"}, &redis.Z{Score: 2, Member: "element2"}).Err()
	if err != nil {
		fmt.Println(err)
	}

	// 修改有序集合元素的分数
	err = client.ZIncrBy(ctx, "mysorted-set", 1, "element1").Err()
	if err != nil {
		fmt.Println(err)
	}

	// 获取有序集合指定范围的元素
	elements, err := client.ZRange(ctx, "mysorted-set", 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("有序集合元素:", elements)
	}

	// 删除有序集合元素
	err = client.ZRem(ctx, "mysorted-set", "element1").Err()
	if err != nil {
		fmt.Println(err)
	}
}

// 事务操作
func transactionOperate(ctx context.Context, client *redis.Client) {
	// 开启事务
	tx := client.TxPipeline()

	// 执行事务操作
	tx.Set(ctx, "key1", "value1", 0)
	tx.Set(ctx, "key2", "value2", 0)

	// 提交事务
	_, err := tx.Exec(ctx)

	if err != nil {
		fmt.Println(err)
	}

	value, err := client.Get(ctx, "key1").Result()
	if err == redis.Nil {
		fmt.Println("key1不存在")
	} else if err != nil {
		fmt.Println("value error:", err)
	} else {
		fmt.Println("key1:", value)
	}
}
