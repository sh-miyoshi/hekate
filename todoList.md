# TODO List

## new commands

- CLI Tool(jwtctl)
  - ServerのAPIをたたくためのコマンドラインツール
- Gateway
  - Backendのユーザープログラムに対してアクセス制御するようなツール
  - keycloak-gatekeeperのようなものを想定
  - tokenを認証するようなAPIを追加(To Server)

## server application enhancement

- open id connect連携
- audit log
- projectのimport/export
- GUI画面の追加
- 設定項目の追加
  - パスワードポリシー
  - refresh tokenのrevoke
  - encrypt_type
- (project/user) enabledの有効化
- SAML対応
- テストの追加
  - ロジック部分のunit test
  - API部分のテスト
- 各種APIの実装
  - user api
- custom roleの有効化
- その他のDB Handlerの実装

## operation enhancement

- add kubernetes yaml file
- write usage to README.md
- create public docker image
- configure CI

## An Idea of revoke refresh token

- add api for `revoke session (useID, sessionID)`
- add api for `get session detail (project, useID, sessionID)`
