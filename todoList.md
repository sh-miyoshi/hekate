# TODO List

## new commands

- Gateway
  - Backendのユーザープログラムに対してアクセス制御するようなツール
  - keycloak-gatekeeperのようなものを想定
- GUI画面(portal)

## jwt-server application enhancement

- 各種APIの実装
  - openid-connect用のAPI
    - /api/v1/project/{projectName}/.well-known/openid-configuration
    - /api/v1/project/{projectName}/protocol/openid-connect/auth
    - /api/v1/project/{projectName}/protocol/openid-connect/token
    - /api/v1/project/{projectName}/protocol/openid-connect/userinfo
    - /api/v1/project/{projectName}/protocol/openid-connect/certs
  - sessionの詳細取得(引数: project, useID, sessionID)
  - 特定ユーザのログアウト(session全削除)
- expireしないclientの登録
- テストの追加
  - ロジック部分のunit test
  - API部分のテスト
- apiのvalidateionの追加
  - user nameなど
- http errorの充実
  - example: [facebook for developers](https://developers.facebook.com/docs/messenger-platform/reference/send-api/error-codes?locale=ja_JP)
- audit log
  - time
  - resource type (or url path and method)
  - client
  - success or failed
- projectのimport/export
- 設定項目の追加
  - パスワードポリシー
  - encrypt_type
- custom roleの有効化
- SAML対応
- (project/user) enabledの有効化
- user federation
  - user情報を外部に保存し、それと連携する
- redirect_urlの設定
- LDAP連携？

## CLI tool(jwtctl) enhancement

- 各APIへの対応
  - project get list
  - project get
  - project update
  - project delete
  - user create
  - user get list
  - user get
  - user update
  - user delete
  - user role add
  - user role delete
- default config pathの修正
- Production向け実行ファイルの作成

## operation enhancement

- add kubernetes yaml file
- write usage to README.md
- create public docker image
- configure CI
