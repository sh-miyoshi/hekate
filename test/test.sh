#!/bin/bash

URL="http://localhost:8080/api/v1"
PROJECT_ID="master"

master_token=`curl -s -X POST -d '@token_request.json' $URL/project/$PROJECT_ID/token | jq -r .accessToken`

# Create New Project
project=`curl -s -X POST -d '@project_create.json' -H "Authorization: Bearer $master_token" $URL/project | jq .`
project_id=`echo $project | jq -r .id`

# Update Project
curl -s -X PUT -d '@project_update.json' -H "Authorization: Bearer $master_token" $URL/project/$project_id

# Get All Projects
curl -s -H "Authorization: Bearer $master_token" $URL/project
