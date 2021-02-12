#!/bin/bash

SERVER_ADDR="http://localhost:18443"
URL="$SERVER_ADDR/userapi/v1"

function test_api() {
	url=$1
	method=$2
	token=$3

	if [ $# = 4 ]; then
		input=$4
		result=`curl --insecure -s -X $method -d "@$input" \
		  -H "Authorization: Bearer $token" $url \
		  -o /dev/null -w '%{http_code}'`
		if [ $result -ge 300 ]; then
			echo "Failed to $method to $url. status code: $result"
			exit 1
		fi
	else
		result=`curl --insecure -s -X $method \
		  -H "Authorization: Bearer $token" $url \
		  -o /dev/null -w '%{http_code}'`
		if [ $result -ge 300 ]; then
			echo "Failed to $method to $url. status code: $result"
			exit 1
		fi
	fi
}

function test_api_return_json() {
	url=$1
	method=$2
	token=$3

	if [ $# = 4 ]; then
		input=$4
		result=`curl --insecure -s -X $method -d "@$input" \
		  -H "Authorization: Bearer $token" $url \
		  | jq .`
		if [ $? != 0 ]; then
			echo "Failed to $method to $url"
			exit 1
		fi
		echo $result
	else
		result=`curl --insecure -s -X $method \
		  -H "Authorization: Bearer $token" $url \
		  | jq .`
		if [ $? != 0 ]; then
			echo "Failed to $method to $url"
			exit 1
		fi
		echo $result
	fi
}

curl $SERVER_ADDR/healthz -s -o /dev/null
if [ $? != 0 ]; then
  echo "Before test, please run a server"
  exit 1
fi

# login
token_info=`curl --insecure -s -X POST $SERVER_ADDR/authapi/v1/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin" \
  -d "password=password" \
  -d "client_id=portal" \
  -d 'grant_type=password'`
token=`echo $token_info | jq -r .access_token`

# get user id
userID=`curl --insecure -s -H "Authorization: Bearer $token" $SERVER_ADDR/adminapi/v1/project/master/user | jq -r .[0].id`
echo "User ID: $userID"

# change-password
test_api "$URL/project/master/user/$userID/change-password" POST $token "inputs/user_change_password.json"

# logout
test_api "$URL/project/master/user/$userID/logout" POST $token
echo "success to user logout"