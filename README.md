# jwt-server

## Overview

jwt-serverはJWT Tokenをできるだけ簡単に取得するためのサーバです。  
主に検証用にお使いください。

## 使い方

TODO

## OS Environment

|環境変数名|内容|Default値|
|---------|----|:-------:|
|JWT_SERVER_ADMIN_NAME|system adminユーザの名前|admin|
|JWT_SERVER_ADMIN_PASSWORD|adminユーザのパスワード|password|
|JWT_SERVER_TOKEN_ISSUER|JWT TokenのIssuer|jwt-server|
|JWT_SERVER_TOKEN_EXPIRED_TIME|JWT Tokenの有効期限(秒)|3600|
