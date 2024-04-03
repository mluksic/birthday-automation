terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.25.0"
    }
  }
}

provider "aws" {
  # Configuration options
  region = "eu-central-1"
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "birthdayAutomation" {
  # If the file is not in the current working directory you will need to include a 
  # path.module in the filename.
  filename      = "../birthdays.zip"
  function_name = "birthdayAutomation"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "birthdays"
  runtime       = "go1.x"

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = filebase64sha256("../birthdays.zip")


  environment {
    variables = {
      FROM_EMAIL      = var.FROM_EMAIL
      APP_PASSWORD    = var.APP_PASSWORD
      EMAIL_RECEIVERS = var.EMAIL_RECEIVERS
      SMTP_HOST       = var.SMTP_HOST
      SMTP_PORT       = var.SMTP_PORT
    }
  }
}


#####################
## EXTRA RESOURCES ##
#####################

# Create cloudwatch event rule
resource "aws_cloudwatch_event_rule" "every_morning" {
  name                = "runBirthdaysLambda"
  description         = "Fires every morning at 6AM UTC"
  schedule_expression = "cron(0 6 * * ? *)"
}

# Create cloudwatch event target
resource "aws_cloudwatch_event_target" "check_every_morning" {
  rule      = "${aws_cloudwatch_event_rule.every_morning.name}"
  target_id = "lambda"
  arn       = "${aws_lambda_function.birthdayAutomation.arn}"
}

# Create cloudwatch event rule
resource "aws_lambda_permission" "allow_cloudwatch_to_call_check_for_birthdays" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.birthdayAutomation.function_name}"
  principal     = "events.amazonaws.com"
  source_arn    = "${aws_cloudwatch_event_rule.every_morning.arn}"
}
