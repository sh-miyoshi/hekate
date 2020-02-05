# [WIP]jwt-server

## Overview

`jwt-server`はGo言語で記述されたOpenID Connectに対応したシンプルな認証・認可サーバです。  
単一バイナリでユーザー管理と認証・認可処理を実装しています。  

## Project Goal

より速く、よりスケールするシンプルな認証・認可サーバ

## 開発環境

- golang v1.12以上

## 使い方

### サーバの起動

```bash
cd cmd/jwt-server
vi config.yaml
  # 必要に応じて修正
go build
./jwt-server
```

### JWT Tokenの取得

```bash
curl --insecure -s -X POST http://localhost:8080/api/v1/project/master/openid-connect/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=admin" \
  -d "password=password" \
  -d "client_id=admin-cli" \
  -d 'grant_type=password'
```

## All APIs

[api/api.html](api/api.html)、もしくは[test/all_api_test.sh](test/all_api_test.sh)を参照してください。

## Road Map

[TODO List](./todoList.md)を参照してください。
