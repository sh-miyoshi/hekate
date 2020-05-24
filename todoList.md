# TODO List

## server application enhancement

- 設定項目の追加
  - パスワードポリシー
- db manager validationの追加
- Custom RoleにDescriptionを追加
  - API docの修正
  - DBの修正
  - Validationの追加
  - APIの修正
- PUT APIの修正
  - フィールドがない or nullの場合は更新しない
- ユーザーパスワードロック
  - Project Infoに追加
  - APIの修正
    - ログインリクエスト失敗時にロックする
    - 強制ロック解除用のAPI
  - API docの修正
- audit log
  - time
  - resource type (or url path and method)
  - client
  - success or failed
- http errorの充実
  - example: [facebook for developers](https://developers.facebook.com/docs/messenger-platform/reference/send-api/error-codes?locale=ja_JP)
- SQL DBの追加
- DBGCの追加
  - Expiredしたsessionなどを一定期間ごとに削除する
- テストの追加
  - unit test
    - pkg/apiclient
    - pkg/apihandler
    - pkg/client
    - pkg/db
    - pkg/oidc
  - 結合テスト
- OpenID Connect部分のエンハンス
  - AuthRequestに他のパラメータを追加
    - display
    - ui_locales
    - id_token_hint
    - login_hint
    - acr_values
  - code認証失敗時、すべてのtokenを無効化
  - subject_types_supportedにpairwiseをサポート
  - RS256以外のSigining Algorithmのサポート
- kong対応
  - URL: [https://konghq.com/](https://konghq.com/)
- (project/user) enabledの有効化
- API responseのtime formatの見直し
- projectのimport/export
- 各種APIの実装
  - sessionの詳細取得(引数: project, userID, sessionID)
- Client Secretに証明書を追加できるようにする
- filterの追加(user, role)
- SAML対応
- user federation
  - user情報を外部に保存し、それと連携する
- User Authentication HTMLの拡充
  - Client IDを表示(optional)
  - Project名を表示
- LDAP連携？
- http headerの追加

## Portal(Admin Console) enhancement

- Userページの作成
  - Redirect先テスト
  - Userアカウント設定ページ
    - Account Setting(/account)
    - Password(/password)
    - Sessions(/sessions)
    - Audit Log(/logs)
- 各ページの作成
  - User
    - Login Sessionの表示
- validationの追加
  - portal側でリクエストを出す前にはじく
- user role更新時のvalidation
  - cluster系とそれ以外を同時に付与しようとした場合警告を出す
- oidc authのstateチェック
- client_idやproject nameをどうするか(変数化)
- middleware処理
  - roleが足りない(masterプロジェクトにいない、cluster操作権限がない?)
    - 強制ログアウト or 白紙のページを見せる(こっちが有力)
- redirect先
  - 前回開いていたページに戻る
- alert画面のcss修正
- headerにユーザー名を表示

- 各ユーザーのアカウント管理画面
  - user名変更
  - パスワード変更

## CLI tool(hctl) enhancement

- いろいろ修正が必要
- configファイルの扱い
  - `no such file or directory`のとき、新規作成
  - permission: 0700, 0600
  - デフォルト値(localhost:18443, master)
- 各APIへの対応
  - project
    - create
      - Allow Grant Typeへの対応
    - update
  - user
    - update
    - role delete
    - session revoke
  - client
    - update
  - customrole
    - get
    - delete
    - update
- outputの修正
  - error出力の方法
  - debug出力
  - error messageの内容
- default config pathの修正
- configコマンドの作成・修正
- Production向け実行ファイルの作成
- support authorization code flow

## new commands

- Gateway
  - Backendのユーザープログラムに対してアクセス制御するようなツール
  - keycloak-gatekeeperのようなものを想定

## operation enhancement

- add kubernetes yaml file
- write usage to README.md
- add release pipeline
  - create public docker image
  - create binary files

## For production

- 初期allowed callback urlの初期化
- CLIのinsecure, timeoutの設定見直し
- Benchmark
