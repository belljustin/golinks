#! /bin/bash

set -e

docker run \
  -p 8080:8080 \
  -v ~/.aws/credentials:/home/nonroot/.aws/credentials:ro \
  -e AWS_PROFILE=golinks \
  -e GOLINKS_STORAGE_TYPE=dynamodb \
  golinks:stable