#!/bin/bash

echo "Start Portal"

echo "Env:"
env | grep HEKATE

npm run build
npm run start

