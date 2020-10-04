# TODO List

## Documents

- 開発者向け構成図

## server application enhancement

- http headerの追加
- db manager validationの追加
- テストの追加
  - unit test
    - pkg/apiclient
    - pkg/apihandler
    - pkg/client
    - pkg/db
    - pkg/oidc
  - 結合テスト
    - DBGCのテスト
- APIの修正
  - Custom RoleにDescriptionを追加
- パスワード以外でのユーザーのログイン
  - 証明書
  - ワンタイムパスワード
  - デバイス認証
- Client Secretに証明書を追加できるようにする
  - portalのアップデートだけでよい？
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
- user federation
  - user情報を外部に保存し、それと連携する
  - LDAP連携？
- SQL DBの追加
- ユーザー情報の追加
  - email
  - first/last name
- kong対応
  - URL: [https://konghq.com/](https://konghq.com/)
- (project/user) enabledの有効化
- projectのimport/export
- User Authentication HTMLの拡充
  - Client IDを表示(optional)
  - Project名を表示
  - user password変更ページの追加
- ~~SAML対応~~

## Portal(Admin Console) enhancement

- 各ページの作成
  - user force password reset
  - audit eventの表示
- ヘルプボタンの追加
  - tooltipなど
- validationの追加
  - portal側でリクエストを出す前にはじく
  - user role更新時のvalidation
    - cluster系とそれ以外を同時に付与しようとした場合警告を出す
- errorの表示の修正
  - エラーページの作成
    - status code: 404, 500
  - black listに登録されているパスワードでユーザーを作成する際のエラーを修正
    - 現在はBad Requestと表示される
- client_idやproject nameをどうするか(変数化)
- redirect先
  - 前回開いていたページに戻る
- alert画面のcss修正
- headerにユーザー名を表示

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
