#!/bin/bash

URL="http://localhost:8080/api/v1"
PROJECT_NAME="master"

token_info=`curl --insecure -s -X POST $URL/project/$PROJECT_NAME/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin" \
  -d "password=password" \
  -d "client_id=admin-cli" \
  -d 'grant_type=password'`
access_token=`echo $token_info | jq -r .access_token`
refresh_token=`echo $token_info | jq -r .refresh_token`

if [ "$access_token" = "null" ]; then
  echo "Failed to get access token"
  exit 1
fi
if [ "$refresh_token" = "null" ]; then
  echo "Failed to get refresh token"
  exit 1
fi
echo "successfully get token"

# Token Update
new_token_info=`curl --insecure -s -X POST $URL/project/$PROJECT_NAME/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "refresh_token=$refresh_token" \
  -d "client_id=admin-cli" \
  -d 'grant_type=refresh_token'`

new_access_token=`echo $new_token_info | jq -r .access_token`
new_refresh_token=`echo $new_token_info | jq -r .refresh_token`

if [ "$new_access_token" = "null" ]; then
  echo "Failed to update access token"
  exit 1
fi
if [ "$new_refresh_token" = "null" ]; then
  echo "Failed to update refresh token"
  exit 1
fi
echo "successfully get new token"

# TODO(Get token by previous refresh token(Expect failed))

# TODO(revoke refresh token, get access token by revoked refresh token)

# TODO(revoke all refresh token in a project)
