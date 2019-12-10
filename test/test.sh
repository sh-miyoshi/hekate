#!/bin/bash

URL="http://localhost:8080/api/v1"
PROJECT_NAME="master"

token_info=`curl -s -X POST -d '@token_request.json' $URL/project/$PROJECT_NAME/token`
master_token=`echo $token_info | jq -r .accessToken`
refresh_token=`echo $token_info | jq -r .refreshToken`

# Create New Project
project=`curl -s -X POST -d '@project_create.json' -H "Authorization: Bearer $master_token" $URL/project | jq .`
project_name=`echo $project | jq -r .name`

# Update Project
curl -s -X PUT -d '@project_update.json' -H "Authorization: Bearer $master_token" $URL/project/$project_name

# Get All Projects
echo "all projects:"
curl -s -H "Authorization: Bearer $master_token" $URL/project

# Token Update
echo "before token: $master_token"
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
echo "after token: $new_token"