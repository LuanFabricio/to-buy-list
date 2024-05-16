package email

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
)

const STMP_GMAIL string = "smtp.gmail.com"
const MIME string = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func SendEmail(to []string, subject string, msg string) error{
	if len(to) <= 0 {
		return errors.New("Empty list of addressee email")
	}

	emailName, found := os.LookupEnv("TBL_EMAIL_NAME")
	if !found {
		return errors.New("E-Mail name not found on env vars")
	}

	email, found := os.LookupEnv("TBL_EMAIL")
	if !found {
		return errors.New("E-Mail address not found on env vars")
	}

	password, found := os.LookupEnv("TBL_PASSWORD")
	if !found {
		return errors.New("E-Mail password not found on env vars")
	}

	auth := smtp.PlainAuth("", email, password, STMP_GMAIL)

	msg_bytes := getMsgBytes(fmt.Sprintf("%s <%s>", emailName, email), to, subject, msg)

	log.Printf("Sending an email to: %s\n", strings.Join(to, ","))
	log.Println(string(msg_bytes))

	err := smtp.SendMail(STMP_GMAIL + ":587", auth, email, to, msg_bytes)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func getMsgBytes(from string, to []string, subject, msg string) []byte {
	msg_bytes := []byte(
		fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %v\r\n", strings.Join(to, ",")) +
		fmt.Sprintf("Subject: %v\r\n", subject) +
		MIME +
		msg,
	)
	return msg_bytes
}

func FetchUsersEmail(db *sql.DB, buyListId int) ([]string, error) {
	const USERS_EMAIL_QUERY = "SELECT u.email FROM buy_list_access bla JOIN users u ON u.id = bla.user_id WHERE  bla.buy_list_id = $1"

	emailsRow, err := db.Query(USERS_EMAIL_QUERY, buyListId)
	if err != nil {
		return nil, err
	}
	defer emailsRow.Close()

	emails := []string{}
	var email string
	for emailsRow.Next() {
		emailsRow.Scan(&email)
		emails = append(emails, email)
	}

	return emails, nil
}
