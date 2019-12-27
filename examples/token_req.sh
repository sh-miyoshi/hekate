#!/bin/bash

url=http://localhost:8080
curl --insecure -s -X POST $url/api/v1/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data-urlencode "username=admin" \
  --data-urlencode "password=password" \
  --data-urlencode 'grant_type=password'