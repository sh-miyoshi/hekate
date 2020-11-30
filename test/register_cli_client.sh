#!/bin/bash

SERVER_ADDR=$1

token=`curl --insecure -s -X POST $SERVER_ADDR/api/v1/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin" \
  -d "password=password" \
  -d "client_id=portal" \
  -d 'grant_type=password' | jq -r .access_token`

status=`curl --insecure -s -X POST -H "Authorization: Bearer $token" \
  "$SERVER_ADDR/api/v1/project/master/client" \
  -d "@inputs/cli_client_create.json" \
  -o /dev/null -w '%{http_code}'`

if [ $status = 409 ] || [ $status = 200 ]; then
  echo "ok."
else
  echo "Failed to create client for cli: $status"
  exit 1
fi
