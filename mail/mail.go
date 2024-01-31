package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

var (
	auth smtp.Auth
	from string
	host string
	port string
	pass string
)

func AuthenticateSMTP() {

	from = os.Getenv("MAIL_FROM")
	pass = os.Getenv("MAIL_PASS")
	host = os.Getenv("MAIL_SMTP_HOST")
	port = os.Getenv("MAIL_SMTP_PORT")

	auth = smtp.PlainAuth("", from, pass, host)

}

func SetupSMTPAuth() {
	AuthenticateSMTP()
}

func SendMail(body string, recipient string, subject string, title string, article string) error {

	message := PrepareMessage(body, recipient, subject, title, article)

	err := smtp.SendMail(host+":"+port, auth, from, []string{recipient}, message)
	if err != nil {
		return err
	}

	fmt.Println("email delivered to: ", recipient)

	return nil

}

func PrepareMessage(body, recipient, subject, title, article string) []byte {

	return []byte("Subject: " + subject + "\r\n" +

		"\r\n" +

		body +

		"\r\n\r\n" +

		"-----------" +

		"\r\n\r\n" +

		title +

		"\r\n" +

		article)

}
