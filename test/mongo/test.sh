#!/bin/bash

# Run mongodb with docker
CONTAINER_NAME="jwt-mongo"
STATUS=`docker ps | grep $CONTAINER_NAME`
if [ "x$STATUS" = "x" ]; then
  docker run --name $CONTAINER_NAME -d -p 27017:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=root \
    -e MONGO_INITDB_ROOT_PASSWORD=example \
    mongo
fi

# Run hekate with mongo driver
cat << EOF > config.yaml
server_port: 18443
server_bind_address: "0.0.0.0"
logfile: ''
debug_mode: true
db:
  type: "mongo"
  connection_string: "mongodb://root:example@localhost:27017"
admin_name: admin
admin_password: password
token_issuer: hekate
token_secret_key: testsecretkey
EOF

cd ../../cmd/server
go build
./server --config=../../test/mongo/config.yaml

# Test API