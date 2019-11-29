# [WIP]jwt-server

## Overview

jwt-serverはJWT Tokenを取得するためのserverです。  
現在絶賛リファクタリング中です。

## 開発環境

- golang v1.12以上

## 使い方

### サーバの起動

```bash
cd cmd/server
vi config.yaml
  # 適当に修正
go build
./server
```

### JWT Tokenの取得

```bash
cat << EOF > token.json
{
    "name": "admin",
    "secret": "password",
    "authType": "password"
}
EOF

curl -X POST -d '@token.json' http://localhost:8080/api/v1/project/master/token
```
