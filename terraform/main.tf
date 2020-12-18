/* Backend */
terraform {
  backend "s3" {
    bucket = "remotebackend"
    key     = "reversevideoapi/terraform.tfstate"
    region  = "us-west-1"
    profile = "jds"
  }
}

/* Providers */
provider "aws" {
  region = var.region
  profile = "jds"
}

/* vars */
variable "region" {
  type = string
  default = "us-west-1"
}

/* remote */
data "terraform_remote_state" "stinkyfingers" {
  backend = "s3"
  config = {
    bucket  = "remotebackend"
    key     = "stinkyfingers/terraform.tfstate"
    region  = "us-west-1"
    profile = "jds"
  }
}


resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_reversevideo_server_lambda"

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

resource "aws_lambda_function" "server_lambda" {
  filename      = "lambda.zip"
  function_name = "reversevideoserverlambda"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "reversevideoapi"
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  runtime = "go1.x"
  timeout = "300"

  environment {
    variables = {
      foo = "bar"
    }
  }
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "reversevideoapi"
  output_path = "lambda.zip"
}

resource "aws_lambda_permission" "server_lambda" {
  statement_id  = "AllowExecutionFromApplicationLoadBalancer"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.server_lambda.arn
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = aws_lb_target_group.server_lambda.arn
}

# IAM
resource "aws_iam_role_policy_attachment" "cloudwatch-attach" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "reverse-video-s3" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.reverse_video_s3.arn
}

resource "aws_iam_policy" "reverse_video_s3" {
  name = "reverse-video-s3"
  description = "gives lambda permissions for s3 bucket"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:*"
      ],
      "Effect": "Allow",
      "Resource": "${aws_lambda_function.server_lambda}"
    }
  ]
}
EOF
}

# ALB
resource "aws_lb_target_group" "server_lambda" {
  name        = "reversevideoapiserverlambda"
  target_type = "lambda"
}

resource "aws_lb_target_group_attachment" "server_lambda" {
  target_group_arn  = aws_lb_target_group.server_lambda.arn
  target_id         = aws_lambda_function.server_lambda.arn
  depends_on        = [aws_lambda_permission.server_lambda]
}

resource "aws_lb_listener_rule" "server_lambda" {
  listener_arn = data.terraform_remote_state.stinkyfingers.outputs.stinkyfingers_https_listener
  priority = 32
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.server_lambda.arn
  }
  condition {
    path_pattern {
      values = ["/reversevideoapi/*"]
    }
  }
  depends_on = [aws_lb_target_group.server_lambda]
}

# S3
resource "aws_s3_bucket" "reversevideo" {
  bucket = "reversevideo"
  acl  "private"
}
