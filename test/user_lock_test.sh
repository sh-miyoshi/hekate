#!/bin/bash

CLI_DIR="../cmd/hctl"

curl localhost:18443/healthz -s -o /dev/null
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

#----------------------------------
# Prepare test project and user
#----------------------------------
echo "building hctl ..."
cd $CLI_DIR
go build
echo "start test"

./hctl login --name admin --password password
./hctl project add --name lock-test --grantTypes password --maxLoginFailure 3 --lockDuration 5 --failureResetTime 10 --userLockEnabled
./hctl user add --project lock-test --name tester --password password

#----------------------------------
# Test 3 times failure
#----------------------------------
test_command login --project lock-test --name tester --password password
test_command_failed login --project lock-test --name tester --password invalid_password
test_command_failed login --project lock-test --name tester --password invalid_password
test_command_failed login --project lock-test --name tester --password invalid_password
test_command_failed login --project lock-test --name tester --password password

#----------------------------------
# Test reset failure state by timeout
#----------------------------------
echo "wait 11[sec] for reset failure state"
sleep 11
test_command_failed login --project lock-test --name tester --password invalid_password
test_command login --project lock-test --name tester --password password

#----------------------------------
# Test force reset by api
#----------------------------------
./hctl login --project lock-test --name tester --password invalid_password
./hctl login --project lock-test --name tester --password invalid_password
./hctl login --project lock-test --name tester --password invalid_password
./hctl login --name admin --password password
test_command user update unlock --project lock-test --name tester
test_command login --project lock-test --name tester --password password

#----------------------------------
# Post-processing
#----------------------------------
./hctl login --name admin --password password
./hctl project delete --name lock-test