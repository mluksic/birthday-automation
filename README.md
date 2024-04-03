# Birthday notifier

Send birthday alerts via email/SMS

## Dependencies

-   [Go](https://go.dev/doc/install)
-   [AWS Lambda](https://aws.amazon.com/lambda/)

## Prerequisites

Download and install:

-   [Go v1.18+](https://go.dev/doc/install)

## Running the app

1. Create `birthdays.csv` file with your birthdays
```
$ cp example_birthdays.csv birthdays.csv
```

2. Run app start command
```bash
$ go run main.go
```

## Build

1. Build binary:

```bash
$ GOOS=linux GOARCH=amd64 go build -o birthdays
```

2. Create ZIP file (binary + CSV):

```bash
$ zip birthdays.zip birthdays birthdays.csv
```

## Deploy

Project uses [Terraform](https://www.terraform.io/) to deploy and provising AWS Lambda function, triggers ect.

Create `secret.tfvars` and fill it with your variables
```
$ cd tf
$ cp example.tfvars secret.tfvars
```

Basic TF commands:
- `terraform plan --var-file="secret.tfvars"` - compares current state and config file, and displays required provision steps
- `terraform apply --var-file="secret.tfvars"` - triggers execution plan (spins up AWS Lambda, creates rules, ect.)

Lambda is currently being triggered every morning at 6AM UTC.

## Test

When Lambda function has been successfully deployed to AWS, run this command:

-   `aws lambda invoke --function-name birthdayAutomation response.json`

## Authors

👤 **Miha Luksic**
