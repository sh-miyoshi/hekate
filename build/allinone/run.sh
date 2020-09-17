#!/bin/bash

function ExecPortal() {
  cd /hekate/portal
  echo "Building portal binaries ..."
  npm run build > portal.log 2>&1
  echo "Start portal"
  echo "Portal log is in /hekate/portal/portal.log"
  npm run start >> portal.log 2>&1
}

if [ "x$SERVER_ADDR" = "x" ]; then
  echo "Please set SERVER_ADDR to os env."
  exit 1
fi

if [ "x$SERVER_PORT" = "x" ]; then
  SERVER_PORT=18443
fi

if [ "x$PORTAL_PORT" = "x" ]; then
  PORTAL_PORT=3000
fi

SERVER_HOST=`echo $SERVER_ADDR | sed -e 's|^[^/]*//||' -e 's|:.*$||' -e 's|/.*$||'`
echo "Run Server Host: $SERVER_HOST"

export HEKATE_PORTAL_HOST=$SERVER_HOST
export HEKATE_PORTAL_PORT=$PORTAL_PORT
export HEKATE_PORTAL_ADDR=https://$SERVER_HOST:$PORTAL_PORT
export HEKATE_SERVER_ADDR=https://$SERVER_HOST:$SERVER_PORT

echo "Env:"
env | grep HEKATE

# Run Portal
ExecPortal &

# Run server
cd /hekate/server
echo "Start server"
echo "Server log is in /hekate/server/server.log"
./hekate-server --config=./config.yaml
