package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type User struct {
	fullName string
	birthday string
}

func main() {
	lambda.Start(SendBirthdayAlert)
}

func SendBirthdayAlert() {
	records, err := readFile("birthdays.csv")

	if err != nil {
		log.Fatal(err)
	}

	var todaysBirthdays = getTodaysBirthdays(records)

	if len(todaysBirthdays) > 0 {
		sendBirthdaySMS(todaysBirthdays)
	}
}

func readFile(filename string) ([][]string, error) {
	file, err := os.Open("birthdays.csv")

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

func getTodaysBirthdays(records [][]string) []string {

	var birthdayPersons []string
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
			birthdayPersons = append(birthdayPersons, user.fullName)
		}
	}

	return birthdayPersons
}

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
