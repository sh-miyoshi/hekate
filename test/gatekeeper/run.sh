#!/bin/bash
trap 'kill $(jobs -p)' EXIT

go build -o test-server
./test-server -logfile=test-server.log &
echo "start test backend app"

cd ../../cmd/jwt-server
go build
./jwt-server &
echo "start jwt-server"

cd ../../test/gatekeeper

# wait server up
sleep 1

URL="http://localhost:8080/api/v1"

# Get Master Token
token_info=`curl --insecure -s -X POST $URL/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin" \
  -d "password=password" \
  -d "client_id=gatekeeper" \
  -d 'grant_type=password'`
master_access_token=`echo $token_info | jq -r .access_token`
# echo $master_access_token

# register client
curl --insecure -s -X POST $URL/project/master/client \
  -H "Authorization: Bearer $master_access_token" \
  -H "Content-Type: application/json" \
  -d "@client.json"

ls keycloak-proxy-linux-amd64 > /dev/null 2>&1
if [ $? != 0 ]; then
  wget https://github.com/keycloak/keycloak-gatekeeper/releases/download/v2.3.0/keycloak-proxy-linux-amd64
  chmod +x keycloak-proxy-linux-amd64
fi
./keycloak-proxy-linux-amd64 --config=config.yml > gatekeeper.log 2>&1 &

# echo "access without gatekeeper"
# curl http://localhost:10000/hello
# echo ""

echo "access with gatekeeper"
curl http://localhost:3000/hello \
  -H "Authorization: Bearer $master_access_token"
echo ""

# wait