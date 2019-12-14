#!/bin/bash

URL="http://localhost:8080/api/v1"

function test_api() {
	url=$1
	method=$2
	token=$3

	if [ $# = 4 ]; then
		input=$4
		result=`curl -s -X $method -d "@$input" \
		  -H "Authorization: Bearer $token" $url \
		  -o /dev/null -w '%{http_code}'`
		if [ $result -lt 300 ]; then
			echo "success"
		else
			echo "failed"
		fi
	else
		result=`curl -s -X $method \
		  -H "Authorization: Bearer $token" $url \
		  -o /dev/null -w '%{http_code}'`
		if [ $result -lt 300 ]; then
			echo "success"
		else
			echo "failed"
		fi
	fi
}

# Get Master Token
token_info=`curl -s -X POST -d '@token_request.json' $URL/project/master/token`
master_access_token=`echo $token_info | jq -r .accessToken`


# Project Create
result=`test_api "$URL/project" POST $master_access_token 'project_create.json'`
echo $result
if [ $result != "success" ]; then
	echo "Failed to create project"
	exit 1
fi

# All Project Get
result=`test_api "$URL/project" GET $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to get project list"
	exit 1
fi

# Project Get
result=`test_api "$URL/project/new-project" GET $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to get project detail"
	exit 1
fi

# Project Update
result=`test_api "$URL/project/new-project" PUT $master_access_token "project_update.json"`
echo $result
if [ $result != "success" ]; then
	echo "Failed to update project"
	exit 1
fi

# Project Delete
result=`test_api "$URL/project/new-project" DELETE $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to delete project"
	exit 1
fi


# User Create
# TODO

# All User Get
result=`test_api "$URL/project/master/user" GET $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to get user list"
	exit 1
fi


# User Get
# TODO

# User Update
# TODO

# User Delete
# TODO

# Append User Role
# TODO

# Remove User Role
# TODO
