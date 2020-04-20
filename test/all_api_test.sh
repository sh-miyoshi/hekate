#!/bin/bash

URL="http://localhost:18443/api/v1"

function test_api() {
	url=$1
	method=$2
	token=$3

	if [ $# = 4 ]; then
		input=$4
		result=`curl --insecure -s -X $method -d "@$input" \
		  -H "Authorization: Bearer $token" $url \
		  -o /dev/null -w '%{http_code}'`
		if [ $result -lt 300 ]; then
			echo "success"
		else
			echo "failed"
		fi
	else
		result=`curl --insecure -s -X $method \
		  -H "Authorization: Bearer $token" $url \
		  -o /dev/null -w '%{http_code}'`
		if [ $result -lt 300 ]; then
			echo "success"
		else
			echo "failed"
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
			exit 1
		fi
		echo $result
	else
		result=`curl --insecure -s -X $method \
		  -H "Authorization: Bearer $token" $url \
		  | jq .`
		if [ $? != 0 ]; then
			exit 1
		fi
		echo $result
	fi
}

# Get Master Token
token_info=`curl --insecure -s -X POST $URL/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin" \
  -d "password=password" \
  -d "client_id=portal" \
  -d 'grant_type=password'`
master_access_token=`echo $token_info | jq -r .access_token`

# Project Create
result=`test_api "$URL/project" POST $master_access_token 'inputs/project_create.json'`
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
result=`test_api "$URL/project/new-project" PUT $master_access_token "inputs/project_update.json"`
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
result=`test_api_return_json "$URL/project/master/user" POST $master_access_token 'inputs/user_create.json'`
if [ $? != 0 ]; then
	echo "Failed to create user"
	exit 1
fi
echo "success"
userID=`echo $result | jq -r .id`

# All User Get
result=`test_api "$URL/project/master/user" GET $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to get user list"
	exit 1
fi

# User Get
result=`test_api "$URL/project/master/user/$userID" GET $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to get user"
	exit 1
fi

# User Update
result=`test_api "$URL/project/master/user/$userID" PUT $master_access_token 'inputs/user_update.json'`
echo $result
if [ $result != "success" ]; then
	echo "Failed to update user"
	exit 1
fi

# Add User Role
result=`test_api "$URL/project/master/user/$userID/role/read-project" POST $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to add role to user"
	exit 1
fi

# Delete User Role
result=`test_api "$URL/project/master/user/$userID/role/read-project" DELETE $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to delete role from user"
	exit 1
fi

# User Password Change
## Get User Token
token_info=`curl --insecure -s -X POST $URL/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=new-user" \
  -d "password=password" \
  -d "client_id=portal" \
  -d 'grant_type=password'`
user_access_token=`echo $token_info | jq -r .access_token`
## Change Password
result=`test_api "$URL/project/master/user/$userID/change-password" POST $user_access_token 'inputs/change-password.json'`
echo $result
if [ $result != "success" ]; then
	echo "Failed to change password"
	exit 1
fi

# User Delete
result=`test_api "$URL/project/master/user/$userID" DELETE $master_access_token 'inputs/user_change_password.json'`
echo $result
if [ $result != "success" ]; then
	echo "Failed to change user password"
	exit 1
fi

# Client Create
result=`test_api "$URL/project/master/client" POST $master_access_token 'inputs/client_create.json'`
echo $result
if [ $result != "success" ]; then
	echo "Failed to create client"
	exit 1
fi
clientID="oidc-client"

# All Client Get
result=`test_api "$URL/project/master/client" GET $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to get client list"
	exit 1
fi

# Client Get
result=`test_api "$URL/project/master/client/$clientID" GET $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to get client"
	exit 1
fi

# Client Update
result=`test_api "$URL/project/master/client/$clientID" PUT $master_access_token 'inputs/client_update.json'`
echo $result
if [ $result != "success" ]; then
	echo "Failed to update client"
	exit 1
fi

# Client Delete
result=`test_api "$URL/project/master/client/$clientID" DELETE $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to delete client"
	exit 1
fi

# Custom Role Create
result=`test_api_return_json "$URL/project/master/role" POST $master_access_token 'inputs/role_create.json'`
if [ $? != 0 ]; then
	echo "Failed to create custom role"
	exit 1
fi
echo "success"
roleID=`echo $result | jq -r .id`

# All Custom Role Get
result=`test_api "$URL/project/master/role" GET $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to get custom role list"
	exit 1
fi

# Custom Role Get
result=`test_api "$URL/project/master/role/$roleID" GET $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to get custom role"
	exit 1
fi

# Custom Role Update
result=`test_api "$URL/project/master/role/$roleID" PUT $master_access_token 'inputs/role_update.json'`
echo $result
if [ $result != "success" ]; then
	echo "Failed to update custom role"
	exit 1
fi

# Custom Role Delete
result=`test_api "$URL/project/master/role/$roleID" DELETE $master_access_token`
echo $result
if [ $result != "success" ]; then
	echo "Failed to delete custom role"
	exit 1
fi
