#!/bin/bash

# Token get
token_info=`curl --insecure -s -X POST http://localhost:18443/adminapi/v1/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin" \
  -d "password=password" \
  -d "client_id=portal" \
  -d 'grant_type=password'`
master_access_token=`echo $token_info | jq -r .access_token`

# Register callback URL
curl --insecure -s -X PUT \
  -d "{\"id\":\"portal\",\"access_type\":\"public\",\"allowed_callback_urls\":[\"http://localhost:3000/callback\"]}" \
  -H "Authorization: Bearer $master_access_token" \
  http://localhost:18443/adminapi/v1/project/master/client/portal
