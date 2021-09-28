package main

import (
	"context"
	"fmt"
	"log"

	redis "github.com/go-redis/redis/v8"
)

func main() {
	err := example()
	if err != nil {
		log.Fatal(err)
	}
}

var ctx = context.Background()

func example() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		return err
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		return err
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		return err
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
	return nil
}
