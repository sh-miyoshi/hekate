#!/bin/bash

function ExecPortal() {
  cd /myapp/portal
  npm run build
  npm run start > portal.log 2>&1
}

if [ "x$SERVER_ADDR" = "x" ]; then
  echo "Please set SERVER_ADDR to os env."
  exit 1
fi

SERVER_HOST=`echo $SERVER_ADDR | sed -e 's|^[^/]*//||' -e 's|:.*$||' -e 's|/.*$||'`
echo "Run Server Host: $SERVER_HOST"

export HEKATE_PORTAL_HOST=$SERVER_HOST
export HEKATE_PORTAL_PORT=3000
# TODO(use https)
export HEKATE_PORTAL_ADDR=http://$SERVER_HOST:3000
export HEKATE_SERVER_ADDR=https://$SERVER_HOST:8080

echo "Env:"
env | grep HEKATE

# Run Portal
ExecPortal &

# Run server
cd /myapp/server
./hekate-server --config=./config.yaml
