package redis

import (
	"context"
	"fmt"
)



func ReadPassword(email string) (string, error) {
	fmt.Println("pass read go")
	ctx := context.Background()
	rdb := RedisConn()
	defer rdb.Close()

	code, err := rdb.Get(ctx, "email:"+email).Result()
	if err != nil {
		fmt.Println(err, "44444444444")
		return "", err
	}

	return code, nil
}
