#!/bin/bash

# set env
export HEKATE_PORTAL_ADDR=http://localhost:3000

# run server
cd ../cmd/hekate
go build
./hekate --config=config.yaml &

# run portal
cd ../portal
npm run dev