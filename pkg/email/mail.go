package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func email() {
	from := os.Getenv("FromEmailAddr")
	password := os.Getenv("SMTPpswd")
	toEmail := os.Getenv("ToEmailAddr")
	to := []string{toEmail}
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	subject := "Subject: Our Golang Email\n"
	body := "our first email!"
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("go check your email")
}
