package email

import (
	"auth_service/configs"
	redis "auth_service/storage/redis"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"time"

	"golang.org/x/exp/rand"
)

func Email(email string) (string, error) {
	client := redis.RedisConn()
	defer client.Close()

	randomNumber := rand.Intn(900000) + 100000
	code := strconv.Itoa(randomNumber)

	err := client.Set(context.Background(), "email:"+email, code, time.Minute*5).Err()
	if err != nil {
		return "", err
	}

	err = client.Set(context.Background(), "e:email", email, time.Minute*5).Err()
	if err != nil {
		return "", err
	}

	err = client.Set(context.Background(), code, email, time.Minute*5).Err()
	if err != nil {
		return "", err
	}

	code1, err := client.Get(context.Background(), "email:"+email).Result()
	if err != nil {
		fmt.Println(err, "88888888888888")
		return "", err
	}
	fmt.Println(code1, "sadfasdfasdfadsfsadfasd")
	err = SendCode(email, code1)

	if err != nil {
		return "", err
	}

	return code1, nil
}

func SendCode(email string, code string) error {
	cfg := configs.Load()
	from := cfg.Email
	password := cfg.Password

	to := []string{
		email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	var body bytes.Buffer

	body.Write([]byte(fmt.Sprintf("Subject: Your verification code \n%s\n\n", code)))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
		return err
	}

	return nil
}
