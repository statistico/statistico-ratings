#!/bin/bash

set -e

aws ecr get-login --no-include-email --region $AWS_DEFAULT_REGION | bash

docker tag "statisticoratings_console" "$AWS_ECR_ACCOUNT_URL/statistico-ratings:$CIRCLE_SHA1"
docker push "$AWS_ECR_ACCOUNT_URL/statistico-ratings:$CIRCLE_SHA1"
