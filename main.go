package main

import (
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type Notifier interface {
	Notify(msg string) error
}

type EmailNotifier struct {
	fromEmail string
	password  string
	smtpHost  string
	smtpPort  string
}

func NewEmailNotifier() *EmailNotifier {
	return &EmailNotifier{
		fromEmail: os.Getenv("FROM_EMAIL"),
		password:  os.Getenv("APP_PASSWORD"),
		smtpHost:  os.Getenv("SMTP_HOST"),
		smtpPort:  os.Getenv("SMTP_PORT"),
	}
}

func (n *EmailNotifier) Notify(msg string) error {
	message := []byte(fmt.Sprintf("Subject: današnji rojstni dnevi \n\n %s\n", msg))

	receivers := strings.Split(os.Getenv("EMAIL_RECEIVERS"), ",")
	auth := smtp.PlainAuth("", n.fromEmail, n.password, n.smtpHost)
	addr := fmt.Sprintf("%s:%s", n.smtpHost, n.smtpPort)

	err := smtp.SendMail(addr, auth, n.fromEmail, receivers, message)
	if err != nil {
		return err
	}

	return nil
}

type User struct {
	fullName string
	birthday string
}

func main() {
	lambda.Start(birthdayAutomation)
}

func birthdayAutomation() error {
	records, err := readFile("birthdays.csv")
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	todayBirthdays := getTodayBirthdays(records)
	notifier := NewEmailNotifier()

	if len(todayBirthdays) > 0 {
		msg := createMsg(todayBirthdays)

		err := notifier.Notify(msg)
		if err != nil {
			log.Printf("%v", err)
			return err
		}

		log.Println("Notification sent successfully")
	}

	return nil
}

func readFile(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Print(err)
		return [][]string{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// skip first line
	if _, err := reader.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func getTodayBirthdays(records [][]string) []User {
	birthdayPersons := make([]User, 0, 10)
	currentDate := time.Now().Local()

	for _, record := range records {
		user := User{
			fullName: record[0],
			birthday: record[1],
		}

		parsedUserBirthday, err := time.Parse("2006-01-02", user.birthday)
		if err != nil {
			fmt.Printf("Could not parse time: %v", err)
		}

		if currentDate.Day() == parsedUserBirthday.Day() && currentDate.Month() == parsedUserBirthday.Month() {
			birthdayPersons = append(birthdayPersons, user)
		}
	}

	return birthdayPersons
}

func createMsg(todayBirthdays []User) string {
	var msg string

	for _, birthday := range todayBirthdays {
		msg += strings.Join([]string{birthday.fullName, birthday.birthday}, " - ")
		msg += "\n"
	}

	return msg
}
