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

# Token Update
cat << EOF > new_token_request.json
{
    "name": "admin",
    "secret": "$refresh_token",
    "authType": "refresh"
}
EOF
new_token_info=`curl -s -X POST -d '@new_token_request.json' $URL/project/$PROJECT_NAME/token`
new_token=`echo $new_token_info | jq -r .accessToken`

# TODO(Get token by previous refresh token(Expect failed))
failed_token=`curl -s -X POST -d '@new_token_request.json' $URL/project/$PROJECT_NAME/token | jq -r .accessToken`
echo $failed_token

rm -f new_token_request.json

# TODO(revoke refresh token)
# TODO(get access token by revoked refresh token)

# TODO(revoke all refresh token in a project)