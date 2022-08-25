# Birthday notifier

Send SMS birthday notification

## Dependencies

-   [Go](https://go.dev/doc/install)
-   [AWS Lambda](https://aws.amazon.com/lambda/)
-   [Twilio](https://www.twilio.com/sms)

### Prerequisites

Download and install:

-   [Go](https://go.dev/doc/install)

### Running the app

1. Create `birthdays.csv` file with people's birthdays
2. Change environment variables
3. Run command below to start the project

```bash
$ go run main.go
```

### Build

1. Build binary:

```bash
$ GOOS=linux GOARCH=amd64 go build -o birthdays
```

2. Create ZIP file (binary + CSV):

```bash
$ zip birthdays.zip birthdays birthdays.csv
```

### Deploy

Project uses [Terraform](https://www.terraform.io/) to deploy and provising AWS Lambda function, triggers ect.

Basic TF commands:

-   `terraform plan` - compares current state and config file, and displays required provision steps
-   `terraform apply` - triggers execution plan (spins up AWS Lambda, creates rules, ect.)

Lambda is currently being triggered every workday (MON - FRI) at 6pm UTC.

### Test

When Lambda function has been successfuly deployed to AWS, run this command:

-   `aws lambda invoke --function-name birthdayAutomation response.json`

## Authors

👤 **Miha Luksic**
