# TODO List

## Documents

- login pageの修正方法
- kubernetes install
  - portalのインストール
- access control
  - 後始末
- ユーザー管理機能
  - 掲示板アプリにログイン機能を実装する方法
- 開発者向け構成図

## server application enhancement

- db manager validationの追加
- DBGC
  - session cleanup
- ログインディレクトリの構造チェック
- masterプロジェクト初期構成時にpassword grantを外す
- SSOの有効期限をconfigファイルで設定できるようにする
- errors.WriteOAuthErrorのエラーハンドリング
- テストの追加
  - unit test
    - pkg/apiclient
    - pkg/apihandler
    - pkg/client
    - pkg/db
    - pkg/db/memory
      - user filter test
    - pkg/oidc
  - 結合テスト
    - DBGCのテスト
    - mongo db test
      - custom role
      - login session
      - session
      - user
    - transaction test
- APIの修正
  - Custom RoleにDescriptionを追加
  - Audit EventにFilter Ruleを追加
    - Project Name, user ID, client ID など
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
  - ID Tokenにほかのパラメータを追加
    - acr
    - amr
    - azp
  - response mode: form_postのサポート
  - code認証失敗時、すべてのtokenを無効化
  - subject_types_supportedにpairwiseをサポート
  - RS256以外のSigining Algorithmのサポート
  - auth requestをparseする
  - type noneのサポート
  - Support PKCE
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
  - user revoke session
- ヘルプボタンの追加
  - tooltipなど
- validationの追加
  - portal側でリクエストを出す前にはじく
  - user role更新時のvalidation
    - cluster系とそれ以外を同時に付与しようとした場合警告を出す
- audit eventの表示の修正
  - 日付で絞れるようにする(from, to)
  - ページネーション
- headerにユーザー名を表示
- API timeoutの変数化
- stateのチェック
- ~~redirect先~~
  - ~~前回開いていたページに戻る~~

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
- Benchmark
