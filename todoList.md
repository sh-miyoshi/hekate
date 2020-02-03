# TODO List

## new commands

- Gateway
  - Backendのユーザープログラムに対してアクセス制御するようなツール
  - keycloak-gatekeeperのようなものを想定
- GUI画面(portal)

## jwt-server application enhancement

- http errorの充実
  - example: [facebook for developers](https://developers.facebook.com/docs/messenger-platform/reference/send-api/error-codes?locale=ja_JP)
  - OIDCのエラーフォーマットに沿う
- 各種APIの実装
  - user role API
    - userに紐付ける
    - token情報に含める
  - openid connect API
    - token revocation
    - implicit flow
    - hybrid flow
  - 特定ユーザのログアウト(session全削除)
  - sessionの詳細取得(引数: project, useID, sessionID)
  - 各リソースのGet APIの見直し
    - 全体検索のみ？queryで検索できるようにする？
- user login session情報をDBに保存する
- テストの追加
  - ロジック部分のunit test
  - API部分のテスト
- https対応
- audit log
  - time
  - resource type (or url path and method)
  - client
  - success or failed
- DBGCの追加
  - Expiredしたsessionなどを一定期間ごとに削除する
- projectのimport/export
- 設定項目の追加
  - パスワードポリシー
  - encrypt_type
- OpenID Connect部分のエンハンス
  - subject_types_supportedにpairwiseをサポート
  - RS256以外のSigining Algorithmのサポート
  - preferred_usernameの追加
- APIのRBACの見直し
- userのパスワード変更のrole見直し
  - 本人のみが変更できるようにする
- SAML対応
- (project/user) enabledの有効化
- user federation
  - user情報を外部に保存し、それと連携する
- redirect_urlの設定
- User Authentication HTMLの拡充
  - Client IDを表示(optional)
  - Project名を表示
- LDAP連携？
- http headerの追加

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
- add release pipeline
  - create public docker image
  - create binary files
