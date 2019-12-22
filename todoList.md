# TODO List

## new commands

- CLI Tool(jwtctl)
  - ServerのAPIをたたくためのコマンドラインツール
- Gateway
  - Backendのユーザープログラムに対してアクセス制御するようなツール
  - keycloak-gatekeeperのようなものを想定
  - tokenを認証するようなAPIを追加(To Server)
- GUI画面(portal)

## server application enhancement

- session dbを外だしする
  - update時に更新されないように
- その他のDB Handlerの実装
  - mongodb driver
    - UserAddRole
    - UserDeleteRole
    - UserNewSession
    - UserRevokeSession
- access tokenの高機能化
  - redirect_url?
- 各種APIの実装
  - sessionの詳細取得(引数: project, useID, sessionID)
  - 特定ユーザのすべてのsessionのrevoke
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
- open id connect連携

## operation enhancement

- add kubernetes yaml file
- write usage to README.md
- create public docker image
- configure CI
