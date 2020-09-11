# TODO List

## Documents

- 構成図

## server application enhancement

- APIの修正
  - Update Project Secret Info APIの追加
- ユーザーパスワードロック
  - APIの追加
    - 強制ロック解除用のAPI
  - API docの修正
- http errorの充実
  - example: [facebook for developers](https://developers.facebook.com/docs/messenger-platform/reference/send-api/error-codes?locale=ja_JP)
  - http.Error関数の置き換え
- ユーザー情報の追加
  - email
  - first/last name
- db manager validationの追加
- Custom RoleにDescriptionを追加
  - API docの修正
  - DBの修正
  - Validationの追加
  - APIの修正
- パスワード以外でのユーザーのログイン
  - 証明書
  - ワンタイムパスワード
  - デバイス認証
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
  - Consentページの追加
    - TokenHandlerからconsent処理
  - AuthRequestに他のパラメータを追加
    - display
    - ui_locales
    - acr_values
  - code認証失敗時、すべてのtokenを無効化
  - subject_types_supportedにpairwiseをサポート
  - RS256以外のSigining Algorithmのサポート
- kong対応
  - URL: [https://konghq.com/](https://konghq.com/)
- (project/user) enabledの有効化
- API responseのtime formatの見直し
- projectのimport/export
- Client Secretに証明書を追加できるようにする
- filterの追加(user, role)
- ~~SAML対応~~
- user federation
  - user情報を外部に保存し、それと連携する
- User Authentication HTMLの拡充
  - Client IDを表示(optional)
  - Project名を表示
- LDAP連携？
- http headerの追加
- user password変更ページの追加

## Portal(Admin Console) enhancement

- 各ページの作成
  - User
    - Login Sessionの表示
    - force password reset
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
- audit eventの表示
- user lock状態の表示
- black listに登録されているパスワードでユーザーを作成する際のエラーを修正
  - 現在はBad Requestと表示される

## CLI tool(hctl) enhancement

- 各APIへの対応
  - user
    - update
    - session revoke
    - password change
  - customrole
    - update
  - audit event
- outputの修正
  - error出力の方法
  - debug出力
  - error messageの内容
- Production向け実行ファイルの作成
- support authorization code flow

## new commands

- Gateway
  - Backendのユーザープログラムに対してアクセス制御するようなツール
  - keycloak-gatekeeperのようなものを想定

## operation enhancement

- add kubernetes yaml file
- add release pipeline
  - create binary files
- login pageの修正方法のドキュメント
- Benchmark
