#!/bin/bash

function ExecPortal() {
  cd /hekate/portal
  npm run build
  npm run start > portal.log 2>&1
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
# TODO(use https)
export HEKATE_PORTAL_ADDR=http://$SERVER_HOST:$PORTAL_PORT
export HEKATE_SERVER_ADDR=https://$SERVER_HOST:$SERVER_PORT

echo "Env:"
env | grep HEKATE

# Run Portal
ExecPortal &

# Run server
cd /hekate/server
./hekate-server --config=./config.yaml
