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

function test_login() {
  project=$1
  name=$2
  passwd=$3

  rawToken=`curl --insecure -s -X POST $SERVER_ADDR/adminapi/v1/project/$project/openid-connect/token \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=$name" \
    -d "password=$passwd" \
    -d "client_id=portal" \
    -d 'grant_type=password'`

  ok=`echo $rawToken | jq .access_token`
  if [ $ok = "null" ]; then
    echo "login with $@ expect success, but got error"
    echo $rawToken
    exit 1
  fi

  # set expires time to tomorrow
  d=`date --rfc-3339=ns -d tomorrow`
  d=`echo $d | sed -e "s/\s/T/g"`

  token=`echo $rawToken | jq ". | { projectName: \"$project\", accessToken: .access_token, accessTokenExpiresDate: \"$d\"}"`
  echo $token > ~/.config/hekate/secret
}

function test_login_failed() {
  project=$1
  name=$2
  passwd=$3

  rawToken=`curl --insecure -s -X POST $SERVER_ADDR/adminapi/v1/project/$project/openid-connect/token \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=$name" \
    -d "password=$passwd" \
    -d "client_id=portal" \
    -d 'grant_type=password'`

  ok=`echo $rawToken | jq .access_token`
  if [ $ok != "null" ]; then
    echo "login with $@ expect failed, but got successed"
    echo $rawToken
    exit 1
  fi
}

#----------------------------------
# Prepare test project and user
#----------------------------------
echo "building hctl ..."
cd $CLI_DIR
go build
echo "start test"

test_login "master" "admin" "password"
./hctl project add --name lock-test --grantTypes password --maxLoginFailure 3 --lockDuration 5 --failureResetTime 10 --userLockEnabled
./hctl user add --project lock-test --name tester --password password

#----------------------------------
# Test 3 times failure
#----------------------------------
test_login "lock-test" "tester" "password"
test_login_failed "lock-test" "tester" "invalid_password"
test_login_failed "lock-test" "tester" "invalid_password"
test_login_failed "lock-test" "tester" "invalid_password"
test_login_failed "lock-test" "tester" "password"

#----------------------------------
# Test reset failure state by timeout
#----------------------------------
echo "wait 11[sec] for reset failure state"
sleep 11
test_login_failed "lock-test" "tester" "invalid_password"
test_login "lock-test" "tester" "password"

#----------------------------------
# Test force reset by api
#----------------------------------
test_login "lock-test" "tester" "password"
test_login_failed "lock-test" "tester" "invalid_password"
test_login_failed "lock-test" "tester" "invalid_password"
test_login_failed "lock-test" "tester" "invalid_password"
test_login "master" "admin" "password"
test_command user update unlock --project lock-test --name tester
test_login "lock-test" "tester" "password"

#----------------------------------
# Post-processing
#----------------------------------
test_login "master" "admin" "password"
./hctl project delete --name lock-test