#!/bin/bash

SERVER_ADDR="http://localhost:18443"
URL="$SERVER_ADDR/api/v1"

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

# Get Master Token
token_info=`curl --insecure -s -X POST $URL/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin" \
  -d "password=password" \
  -d "client_id=portal" \
  -d 'grant_type=password'`
master_access_token=`echo $token_info | jq -r .access_token`

# Project Create
test_api "$URL/project" POST $master_access_token 'inputs/project_create.json'
echo "success to project create"

# All Project Get
test_api "$URL/project" GET $master_access_token
echo "success to all project get"

# Project Get
test_api "$URL/project/new-project" GET $master_access_token
echo "success to project get"

# Project Update
test_api "$URL/project/new-project" PUT $master_access_token "inputs/project_update.json"
echo "success to project update"

# Keys Get
test_api "$URL/project/new-project/keys" GET $master_access_token
echo "success to get project secret"

# Keys Reset
test_api "$URL/project/new-project/keys/reset" POST $master_access_token
echo "success to reset project secret"

# Project Delete
test_api "$URL/project/new-project" DELETE $master_access_token
echo "success to project delete"

# Custom Role Create
result=`test_api_return_json "$URL/project/master/role" POST $master_access_token 'inputs/role_create.json'`
echo "success to custom role create"
roleID=`echo $result | jq -r .id`

# All Custom Role Get
test_api "$URL/project/master/role" GET $master_access_token
echo "success to all custom role get"

# Custom Role Get
test_api "$URL/project/master/role/$roleID" GET $master_access_token
echo "success to custom role get"

# Custom Role Update
test_api "$URL/project/master/role/$roleID" PUT $master_access_token 'inputs/role_update.json'
echo "success to custom role update"

# User Create
result=`test_api_return_json "$URL/project/master/user" POST $master_access_token 'inputs/user_create.json'`
echo "success to user create"
userID=`echo $result | jq -r .id`

# Get User Token
token_info=`curl --insecure -s -X POST $URL/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=user1" \
  -d "password=password" \
  -d "client_id=portal" \
  -d 'grant_type=password'`
user_access_token=`echo $token_info | jq -r .access_token`

# All User Get
test_api "$URL/project/master/user" GET $master_access_token
echo "success to all user get"

# User Get
result=`test_api_return_json "$URL/project/master/user/$userID" GET $master_access_token`
echo "success to user get"
sessionID=`echo $result | jq -r .sessions[0]`

# User Update
test_api "$URL/project/master/user/$userID" PUT $master_access_token 'inputs/user_update.json'
echo "success to user update"

# Add User Role
test_api "$URL/project/master/user/$userID/role/read-project" POST $master_access_token
echo "success to add user role"

# Delete User Role
test_api "$URL/project/master/user/$userID/role/read-project" DELETE $master_access_token
echo "success to delete user role"

# Add Custom Role to User
test_api "$URL/project/master/user/$userID/role/$roleID" POST $master_access_token

# User Password Change
test_api "$URL/project/master/user/$userID/reset-password" POST $master_access_token 'inputs/change-password.json'
echo "success to reset password"

# Get Session
test_api "$URL/project/master/session/$sessionID" GET $master_access_token
echo "success to get session"

# Delete Session
test_api "$URL/project/master/session/$sessionID" DELETE $master_access_token
echo "success to delete session"

# User Delete
test_api "$URL/project/master/user/$userID" DELETE $master_access_token 'inputs/user_change_password.json'
echo "success to user delete"

# Custom Role Delete
test_api "$URL/project/master/role/$roleID" DELETE $master_access_token
echo "success to custom role delete"

# Client Create
test_api "$URL/project/master/client" POST $master_access_token 'inputs/client_create.json'
echo "success to client create"
clientID="oidc-client"

# All Client Get
test_api "$URL/project/master/client" GET $master_access_token
echo "success to all client get"

# Client Get
test_api "$URL/project/master/client/$clientID" GET $master_access_token
echo "success to client get"

# Client Update
test_api "$URL/project/master/client/$clientID" PUT $master_access_token 'inputs/client_update.json'
echo "success to client update"

# Client Delete
test_api "$URL/project/master/client/$clientID" DELETE $master_access_token
echo "success to client delete"

# Audit Events Get
test_api "$URL/project/master/audit" GET $master_access_token
echo "success to get audit events"

# User Logout
test_api "$URL/project/master/user/$userID/logout" POST $master_access_token
echo "success to user logout"
