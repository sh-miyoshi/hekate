@echo off

rem run server
cd ../cmd/hekate
echo "building server binary ..."
go build

echo "start server in background"
start hekate.exe

rem run portal
cd ../portal
echo "start portal"
npm run dev