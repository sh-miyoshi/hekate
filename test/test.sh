#!/bin/bash

URL="http://localhost:8080/api/v1"
PROJECT_ID="master"

access_token=`curl -X POST -d '@token_request.json' $URL/project/$PROJECT_ID/token | jq -r .accessToken`

# Get All Projects
curl -H "Authorization: Bearer $access_token" $URL/project