#!/bin/bash

CLI_DIR="../cmd/hctl"

curl localhost:8080/healthz -s -o /dev/null
if [ $? != 0 ]; then
  echo "Before test, please run a server"
  exit 1
fi

#----------------------------------
# Prepare
#----------------------------------
cd $CLI_DIR
go build

./hctl login --name admin --password password
./hctl project add --name rbac-test
./hctl user add --project rbac-test --name viewer --password password --systemRoles "read-client"
./hctl user add --project rbac-test --name editor --password password --systemRoles "read-client,write-client"
./hctl project add --name rbac-test-2

#----------------------------------
# Test read/write role
#----------------------------------

# TODO

#----------------------------------
# Test cluster role
#----------------------------------

# TODO

#----------------------------------
# Test access to other project
#----------------------------------

# TODO

echo "Successfully finished"