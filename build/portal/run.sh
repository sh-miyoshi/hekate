#!/bin/bash

cd /hekate/portal
echo "building portal ..."
npm run build

echo "successfully build portal"
echo "start portal ..."
npm run start
