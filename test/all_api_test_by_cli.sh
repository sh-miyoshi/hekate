#!/bin/bash

SERVER_ADDR="http://localhost:18443"
CLI_DIR="../cmd/hctl"
TEST_DIR=$PWD

function test_command() {
  ./hctl $@
  if [ $? != 0 ]; then
    echo "hctl $@ command expect success, but got error"
    exit 1
  fi
}

curl $SERVER_ADDR/healthz -s -o /dev/null
if [ $? != 0 ]; then
  echo "Before test, please run a server"
  exit 1
fi

cd $CLI_DIR
go build

# login
test_command login --project master --name admin --password password

# config
## TODO(get)
## TODO(set)

# project
## create
test_command project add --file $TEST_DIR/inputs/project_create.json

## get all
## get
## update
## delete

# user
## create
## get all
## get
## update
## add user role
## delete user role
## password change
## delete

# client
## create
## get all
## get
## update
## delete

# role
## create
## get all
## get
## update
## delete

# logout
test_command logout

echo "successfully finished test"