package redis

import (
	"context"
	"fmt"
)



func ReadEmail() (string, error) {
	ctx := context.Background()
	fmt.Println(2)
	rdb := RedisConn()
	defer rdb.Close()
	email, err := rdb.Get(ctx, "e:email").Result()
	if err != nil {
		fmt.Println(err, "mana xatolik")
		return "", err
	}
	fmt.Println(email, "1111111111111111111111111")
	return email, nil
}
