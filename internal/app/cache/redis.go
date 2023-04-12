package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"

	"github.com/hablof/omp-bot/internal/config"
	"github.com/hablof/omp-bot/internal/model/logistic"
	"github.com/hablof/omp-bot/internal/service/logistic/mypackage"
)

var _ mypackage.CacheDict = &Cache{}

type Cache struct {
	rc *redis.Client
}

func NewCache(cfg config.Config) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sc := rdb.Ping(ctx)
	if err := sc.Err(); err != nil {
		return nil, err
	}

	rdb = rdb.WithTimeout(100 * time.Millisecond)

	return &Cache{
		rc: rdb,
	}, nil
}

func (c *Cache) SetDescription(unit logistic.Package) error {
	ctx := context.Background()
	keyString := "package" + strconv.FormatUint(unit.ID, 10)

	b, err := json.Marshal(unit)
	if err != nil {
		return err
	}

	if err := c.rc.Set(ctx, keyString, b, 0).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Cache) ReadDescription(id uint64) (*logistic.Package, error) {
	ctx := context.Background()
	keyString := "package" + strconv.FormatUint(id, 10)

	s, err := c.rc.Get(ctx, keyString).Result()
	if err != nil {
		return nil, err
	}

	unit := logistic.Package{}
	if err := json.Unmarshal([]byte(s), &unit); err != nil {
		c.rc.Del(ctx, keyString)
		return nil, err
	}

	return &unit, nil
}

func (c *Cache) RemoveDescription(id uint64) error {
	ctx := context.Background()
	keyString := "package" + strconv.FormatUint(id, 10)

	c.rc.Del(ctx, keyString)

	return nil
}

// err := rdb.Set(ctx, "key", "value", 0).Err()
// if err != nil {
// 	panic(err)
// }

// val, err := rdb.Get(ctx, "key").Result()
// if err != nil {
// 	panic(err)
// }
// fmt.Println("key", val)

// val2, err := rdb.Get(ctx, "key2").Result()
// if err == redis.Nil {
// 	fmt.Println("key2 does not exist")
// } else if err != nil {
// 	panic(err)
// } else {
// 	fmt.Println("key2", val2)
// }
// Output: key value
// key2 does not exist
