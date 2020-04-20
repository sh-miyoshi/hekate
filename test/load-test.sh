#!/bin/bash

if [ $# != 2 ]; then
  echo "Usage $0 <thread-num> <loop-num>"
  exit 1
fi

THREADS=$1
LOOPS=$2

URL="http://localhost:18443/api/v1"
PROJECT_NAME="master"

function test() {
  token_info=`curl --insecure -s -X POST $URL/project/master/openid-connect/token \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=admin" \
    -d "password=password" \
    -d "client_id=portal" \
    -d 'grant_type=password'`
  token=`echo $token_info | jq -r .access_token`

  # user create
  result=`curl -s -X POST -H "Authorization: Bearer $token" \
    "$URL/project/master/user" \
    -d "{\"name\": \"user$1\", \"password\": \"password\"}"`
  userID=`echo $result | jq -r .id`
  # user get
  result=`curl -s -H "Authorization: Bearer $token" "$URL/project/master/user/$userID" -o /dev/null -w '%{http_code}'`
  if [ $result -gt 300 ]; then
    echo "Failed to get user $1"
    exit 1
  fi
  # user delete
  result=`curl -s -X DELETE -H "Authorization: Bearer $token" "$URL/project/master/user/$userID" -o /dev/null -w '%{http_code}'`
  if [ $result -gt 300 ]; then
    echo "Failed to delete user $1"
    exit 1
  fi
}

for i in `seq $LOOPS`; do
  for j in `seq $THREADS`; do
    test $j &
  done
  wait
  echo "Done $i/$LOOPS"
done