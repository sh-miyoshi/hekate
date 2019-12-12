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
    curl -s -X POST -d '@token_request.json' $URL/project/$PROJECT_NAME/token -o /dev/null &
  done
  wait
  echo "Done $i/$LOOPS"
done