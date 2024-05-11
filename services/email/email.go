package email

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
)

const STMP_GMAIL string = "smtp.gmail.com"

func SendEmail(to []string, subject string, msg string) error{
	email, found := os.LookupEnv("TBL_EMAIL")
	if !found {
		return errors.New("E-Mail address not found on env vars")
	}

	password, found := os.LookupEnv("TBL_PASSWORD")
	if !found {
		return errors.New("E-Mail password not found on env vars")
	}

	auth := smtp.PlainAuth("", email, password, STMP_GMAIL)

	msg_bytes := []byte(
		fmt.Sprintf("To: %v\r\n", to) +
		fmt.Sprintf("Subject: %v\r\n", subject) +
		msg,
	)

	err := smtp.SendMail(STMP_GMAIL + ":587", auth, email, to, msg_bytes)

	if err != nil {
		return err
	}

	return nil
}
