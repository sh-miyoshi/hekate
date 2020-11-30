#!/bin/bash

CLI_DIR="../cmd/hctl"
SERVER_ADDR="http://localhost:18443"

curl $SERVER_ADDR/healthz -s -o /dev/null
if [ $? != 0 ]; then
  echo "Before test, please run a server"
  exit 1
fi

function test_command() {
  ./hctl $@
  if [ $? != 0 ]; then
    echo "hctl $@ command expect success, but got error"
    exit 1
  fi
}

function test_command_failed() {
  ./hctl $@
  if [ $? = 0 ]; then
    echo "hctl $@ command expect failed, but successed"
    exit 1
  fi
}

function login() {
  project=$1
  name=$2
  passwd=$3

  rawToken=`curl --insecure -s -X POST $SERVER_ADDR/api/v1/project/$project/openid-connect/token \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=$name" \
    -d "password=$passwd" \
    -d "client_id=portal" \
    -d 'grant_type=password'`

  # set expires time to tomorrow
  d=`date --rfc-3339=ns -d tomorrow`
  d=`echo $d | sed -e "s/\s/T/g"`

  token=`echo $rawToken | jq ". | { projectName: \"$project\", accessToken: .access_token, accessTokenExpiresDate: \"$d\"}"`
  echo $token > ~/.config/hekate/secret
}

#----------------------------------
# Prepare and Test cluster role
#----------------------------------
cd $CLI_DIR
go build

login "master" "admin" "password"
./hctl project add --name rbac-test --grantTypes password
./hctl user add --project rbac-test --name viewer --password password --systemRoles "read-project"
./hctl user add --project rbac-test --name editor --password password --systemRoles "read-project,write-project"
./hctl project add --name rbac-test-2 --grantTypes password
./hctl client add --project rbac-test-2 --id test-client-2 --accessType public

#----------------------------------
# Test read/write role
#----------------------------------

login "rbac-test" "editor" "password"
test_command client add --project rbac-test --id test-client --accessType public
test_command client get --project rbac-test

login "rbac-test" "viewer" "password"
test_command_failed client add --project rbac-test --id test-client-2 --accessType public
test_command client get --project rbac-test

#----------------------------------
# Test access to other project
#----------------------------------

test_command_failed client get --project rbac-test-2

echo "Successfully finished"