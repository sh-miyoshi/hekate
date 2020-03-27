#!/bin/bash

if [ "x$SERVER_ADDR" = "x" ]; then
  echo "Please set SERVER_ADDR to os env."
  exit 1
fi

export HEKATE_PORTAL_ADDR=$SERVER_ADDR

# Run Portal
cd /myapp/portal
npm run start > portal.log 2>&1 &

# Run server
cd /myapp/server
./hekate-server --config=./config.yaml
