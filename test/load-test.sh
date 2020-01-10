#!/bin/bash

if [ $# != 2 ]; then
  echo "Usage $0 <thread-num> <loop-num>"
  exit 1
fi

THREADS=$1
LOOPS=$2

URL="http://localhost:8080/api/v1"
PROJECT_NAME="master"

for i in `seq $LOOPS`; do
  for j in `seq $THREADS`; do
    curl --insecure -s -X POST $URL/project/$PROJECT_NAME/openid-connect/token -o /dev/null \
      -H "Content-Type: application/x-www-form-urlencoded" \
      -d "username=admin" \
      -d "password=password" \
      -d 'grant_type=password' &
  done
  wait
  echo "Done $i/$LOOPS"
done