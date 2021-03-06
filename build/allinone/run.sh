#!/bin/bash

function ExecPortal() {
  port=$1
  cd /hekate/portal
  echo "Run building page server"
  ./portal-building -port $port -https-cert-file="/hekate/secret/tls.crt" -https-key-file="/hekate/secret/tls.key" &
  pid=$!

  echo "Build portal binaries"
  echo "Please wait for a while ..."
  npm run build > portal.log 2>&1

  echo "Stop building page server"
  kill $pid

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

export HEKATE_SERVER_ADDR=https://$SERVER_HOST:$SERVER_PORT
export HEKATE_PORTAL_ADDR=https://$SERVER_HOST:$PORTAL_PORT

cat << EOF > /hekate/portal/.env
HEKATE_SERVER_ADDR = ${HEKATE_SERVER_ADDR}
HEKATE_PORTAL_ADDR = ${HEKATE_PORTAL_ADDR}
HEKATE_PORTAL_PORT = ${PORTAL_PORT}
EOF

echo "Env:"
env | grep HEKATE

# Run mongodb
cd /hekate/mongo
echo "Start mongodb"
echo "Mongodb log is in /hekate/mongo/mongo.log"
mkdir -p data
./mongod --logpath=mongo.log --dbpath="./data" &
cd -

# Run portal
ExecPortal $PORTAL_PORT &

# Run server
cd /hekate/server
echo "Start server"
echo "Server log is in /hekate/server/server.log"
./hekate-server --config=./config.yaml
