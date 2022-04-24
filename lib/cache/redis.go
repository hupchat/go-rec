package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-rec/helpers"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	client *redis.Client = nil
)

type XType string

func use(c context.Context) {
	if client != nil {
		return
	}

	env := helpers.VerifyEnv([]string{
		"REDISCLOUD_URL",
	})
	if env != nil {
		log.Fatal(env)
	}

	cn := strings.Split(os.Getenv("REDISCLOUD_URL"), "redis://")
	//format: redis://user:pass@host:port
	addr := strings.Split(cn[1], "@")[1]
	pass := strings.Split(strings.Split(cn[1], "@")[0], ":")[1]

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	_, err := client.Ping(c).Result()
	if err != nil {
		log.Fatal(err)
	}
}

func ContextProducer(ctx context.Context, stream string, resource map[string]string, message []byte) (string, error) {

	r := map[string]interface{}{
		"timestamp": time.Now().UnixNano(),
		"data":      message,
	}

	for k, v := range resource {
		r[k] = v
	}

	use(ctx)
	return client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: r,
	},
	).Result()
}

func SetHash(ctx context.Context, wg *sync.WaitGroup, clientID, hash, value string) (bool, error) {
	use(ctx)
	r, err := client.HSetNX(ctx, clientID, hash, value).Result()
	wg.Done()
	return r, err
}
