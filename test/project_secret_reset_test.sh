#!/bin/bash

SERVER_ADDR="http://localhost:18443"
URL="$SERVER_ADDR/api/v1"

curl $SERVER_ADDR/healthz -s -o /dev/null
if [ $? != 0 ]; then
  echo "Before test, please run a server"
  exit 1
fi

# Get Master Token
token_info=`curl --insecure -s -X POST $URL/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin" \
  -d "password=password" \
  -d "client_id=portal" \
  -d 'grant_type=password'`
master_access_token=`echo $token_info | jq -r .access_token`

# Get Public Key Info
prev=`curl -k -s -X GET -H "Authorization: Bearer $master_access_token" "$URL/project/master/openid-connect/certs" | jq -r .keys[0].n`

# # Project Secret Reset
curl -k -s -X POST -H "Authorization: Bearer $master_access_token" "$URL/project/master/keys/reset"

# Get Public Key Info
current=`curl -k -s -X GET -H "Authorization: Bearer $master_access_token" "$URL/project/master/openid-connect/certs" | jq -r .keys[0].n`

if [ $prev = $current ]; then
  echo "Failed to reset secret"
  exit 1
fi

echo "Success to reset secret"