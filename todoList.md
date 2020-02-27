# TODO List

## Bugs(?)

- adminユーザーが消せる
- admin-cliも多分消せる
- userが削除された場合のaccess tokenが無効化されない(この仕様はok?)
- redirect_uril登録してなくてもredirectされてしまう？

## new commands

- Gateway
  - Backendのユーザープログラムに対してアクセス制御するようなツール
  - keycloak-gatekeeperのようなものを想定

## jwt-server application enhancement

- GET LIST APIの修正
  - 残り: client, role
  - すべての検索結果を取得するようにする
    - http queryでfilterできるようにする
      - API handlerの修正
      - DB queryの修正
  - API docの修正
    - 残り: client, role, project, user
- roleの割り当てのvalidation
  - write権限のみはだめ(同リソースのreadは必須)
    - 作成・削除時にvalidationをかける
  - masterプロジェクトユーザはcluster-read必須？
  - その他はcluster系は付けられないようにする？
- DBGCの追加
  - Expiredしたsessionなどを一定期間ごとに削除する
- audit log
  - time
  - resource type (or url path and method)
  - client
  - success or failed
- OIDC Authorization Code APIの修正
  - ログインページを返す際に適切なheaderをつける
- 各種APIの実装
  - openid connect API
    - implicit flow
    - hybrid flow
  - sessionの詳細取得(引数: project, useID, sessionID)
- テストの追加
  - ロジック部分のunit test
  - API部分のテスト
- 設定項目の追加
  - パスワードポリシー
  - encrypt_type(signing_method)
- authroztion code flowのユーザーログインの修正
  - ログイン失敗時にエラーを表示する
    - Invalid user name or password
- http errorの充実
  - example: [facebook for developers](https://developers.facebook.com/docs/messenger-platform/reference/send-api/error-codes?locale=ja_JP)
- projectのimport/export
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

## Portal(Admin Console) enhancement

- middleware処理
  - roleが足りない(masterプロジェクトにいない、cluster操作権限がない?)
    - 強制ログアウト or 白紙のページを見せる(こっちが有力)
- redirect先
  - 前回開いていたページに戻る
- alert画面のcss修正
- project選択を一覧に
- 各ページの作成
  - TODO

## CLI tool(jwtctl) enhancement

- 各APIへの対応
  - project
    - update
  - user
    - create
      - file flagの追加
    - update
    - role add
    - role delete
  - client
    - create
    - get
    - delete
    - update
  - customrole
    - create
    - get
    - delete
    - update
- default config pathの修正
- configコマンドの作成・修正
- Production向け実行ファイルの作成

## operation enhancement

- add kubernetes yaml file
- write usage to README.md
- add release pipeline
  - create public docker image
  - create binary files
