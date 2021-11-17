package auth

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func CreateEmailVerHash() (string, string, error) {
	emailVerRandRune := make([]rune, 64)
	for i := 0; i < 64; i++ {
		emailVerRandRune[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes)-1)]
	}
	emailVerPassword := string(emailVerRandRune)
	var emailVerPWhash []byte
	emailVerPWhash, err := bcrypt.GenerateFromPassword([]byte(emailVerPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return string(emailVerPassword), string(emailVerPWhash), nil
}

func CreateTimeout() time.Time {
	currentTime := time.Now()
	return currentTime.Add(time.Minute * 45)
}

func Email(toEmail, body string) {
	fmt.Println("emailemail", toEmail)
	from := os.Getenv("FromEmailAddr")
	password := os.Getenv("SMTPpswd")
	to := []string{toEmail}
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	subject := "Subject: password recovery\n"
	messageBody := "<body><a rel=\"nofollow noopener noreferrer\" target=\"_blank\" href=\"" + body + "></body>"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(subject + mime + messageBody)
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
}
