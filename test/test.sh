#!/bin/bash

URL="http://localhost:8080/api/v1"
PROJECT_NAME="master"

token_info=`curl -s -X POST -d '@token_request.json' $URL/project/$PROJECT_NAME/token`
master_token=`echo $token_info | jq -r .accessToken`
refresh_token=`echo $token_info | jq -r .refreshToken`

if [ "x$master_token" = "x" -o "x$refresh_token" = "x" ]; then
    echo "Failed to get token"
    exit 1
fi

# Get All Users
all_users=`curl -s -H "Authorization: Bearer $master_token" $URL/project/master/user`
if [ "x$all_users" = "x" ]; then
    echo "Failed to get users in master project"
    exit 1
fi

echo "All Users: $all_users"

# Create New Project
project=`curl -s -X POST -d '@project_create.json' -H "Authorization: Bearer $master_token" $URL/project | jq .`
project_name=`echo $project | jq -r .name`

if [ "x$project_name" = "x" ]; then
    echo "Failed to create new project"
    exit 1
fi

# Update Project
update_status=`curl -s -X PUT -d '@project_update.json' -H "Authorization: Bearer $master_token" $URL/project/$project_name -o /dev/null -w '%{http_code}'`

if [ "$update_status" != "204" ]; then
    echo "Failed to update project"
    exit 1
fi

# Get All Projects
all_projects=`curl -s -H "Authorization: Bearer $master_token" $URL/project`

if [ "x$all_projects" = "x" ]; then
    echo "Failed to get projects"
    exit 1
fi

echo "All Projects: $all_projects"

# Token Update
cat << EOF > new_token_request.json
{
    "name": "admin",
    "secret": "$refresh_token",
    "authType": "refresh"
}
EOF
new_token_info=`curl -s -X POST -d '@new_token_request.json' $URL/project/$PROJECT_NAME/token`
rm -f new_token_request.json
new_token=`echo $new_token_info | jq -r .accessToken`

if [ "$master_token" = "$new_token" ]; then
    echo "Failed to update token"
    exit 1
fi

echo "before token: $master_token"
echo "after token: $new_token"
