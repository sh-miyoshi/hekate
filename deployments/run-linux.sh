#!/bin/bash

function stopserver() {
    kill $pid
    exit 0
}

trap stopserver int

# set env
export HEKATE_PORTAL_ADDR=http://localhost:3000

# run server
cd ../cmd/hekate
go build
./hekate --config=config.yaml &
pid=$!

# run portal
cd ../portal
npm run dev