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
# Prepare and Test cluster role
#----------------------------------
cd $CLI_DIR
go build

./hctl login --name admin --password password
./hctl project add --name rbac-test
./hctl user add --project rbac-test --name viewer --password password --systemRoles "read-client"
./hctl user add --project rbac-test --name editor --password password --systemRoles "read-client,write-client"
./hctl project add --name rbac-test-2
./hctl client add --project rbac-test-2 --id test-client-2 --accessType public

#----------------------------------
# Test read/write role
#----------------------------------

./htcl login --name editor --password password
test_command client add --project rbac-test --id test-client --accessType public
test_command client get --project rbac-test

./htcl login --name viewer --password password
test_command_failed client add --project rbac-test --id test-client-2 --accessType public
test_command client get --project rbac-test

#----------------------------------
# Test access to other project
#----------------------------------

test_command_failed client get --project rbac-test-2

echo "Successfully finished"