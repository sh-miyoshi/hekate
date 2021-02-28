@echo off

rem set env
set HEKATE_PORTAL_ADDR=http://localhost:3000

rem run server
cd ../cmd/hekate
go build
start /b /min hekate.exe --config=config.yaml

rem run portal
cd ../portal
npm run dev