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

type User struct {
	fullName string
	birthday string
}

func main() {
	lambda.Start(SendBirthdayAlert)
}

func sendEmail(todayBirthdays []User) {
	from := os.Getenv("fromEmail")
	password := os.Getenv("appPassword")
	receivers := strings.Split(os.Getenv("emailReceivers"), ",")

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "današnji rojstni dnevi"

	var msg string
	for _, birthday := range todayBirthdays {
		msg += strings.Join([]string{birthday.fullName, birthday.birthday}, " - ")
		msg += "\n"
	}

	message := []byte(fmt.Sprintf("Subject: %s \n\n %s\n", subject, msg))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, receivers, message)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Email Sent Successfully!")
}

func SendBirthdayAlert() {
	records, err := readFile("birthdays.csv")

	if err != nil {
		log.Fatal(err)
	}

	var todayBirthdays = getTodayBirthdays(records)

	if len(todayBirthdays) > 0 {
		sendEmail(todayBirthdays)
	}
}

func readFile(filename string) ([][]string, error) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
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
	var birthdayPersons []User
	currentDate := time.Now().Local()

	for _, record := range records {
		user := User{
			fullName: record[0],
			birthday: record[1],
		}

		parsedUserBirthday, err := time.Parse("2006-01-02", user.birthday)

		if err != nil {
			fmt.Println("Could not parse time:", err)
		}

		if currentDate.Day() == parsedUserBirthday.Day() && currentDate.Month() == parsedUserBirthday.Month() {
			birthdayPersons = append(birthdayPersons, user)
		}
	}

	return birthdayPersons
}

/*
func sendBirthdaySMS(birthdays []string) {
	accountSid := os.Getenv("accountSid")
	authToken := os.Getenv("authToken")
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	birthdaysString := strings.Join(birthdays, ", ")

	params := &openapi.CreateMessageParams{}
	params.SetTo(os.Getenv("toPhoneNumber"))
	params.SetBody("Todays birthdays: " + birthdaysString)
	params.SetMessagingServiceSid(os.Getenv("messagingServiceSid"))

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
}
*/
