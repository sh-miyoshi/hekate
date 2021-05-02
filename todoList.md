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

## Server application enhancement

- errors.WriteOAuthErrorのエラーハンドリング
- APIの戻り値のJSONの型名のチェック
- model.LoginSessionの修正
  - time.Time型をやめ、expiresIn int64型にする
- 結合テストの追加
  - DBGCのテスト
  - mongo db test
    - custom role
    - login session
    - session
    - user
  - transaction test
  - PKCE
    - 正常系
    - challengeなしの正常系
    - challnge失敗
- db manager validationの追加
- unit testの追加
  - pkg/apiclient
  - pkg/apihandler
  - pkg/client
  - pkg/db
  - pkg/db/memory
    - user filter test
  - pkg/oidc
- Scopeの修正
  - profileの追加(usernameは設定された時のみ)
  - email_verifiedの追加
- masterプロジェクト初期構成時にpassword grantを外す
- User portalを別に分ける
- ユーザー情報の追加
  - first/last name
- User Authentication HTMLの拡充
  - Client IDを表示(optional)
  - Project名を表示
- kong対応
  - URL: [https://konghq.com/](https://konghq.com/)
- (project/user) enabledの有効化
- APIの修正
  - Custom RoleにDescriptionを追加
  - Audit EventにFilter Ruleを追加
    - Project Name, user ID, client ID など
- SQL DBの追加
- user federation
  - user情報を外部に保存し、それと連携する
  - LDAP連携？
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
- TOTPで前後1つも許可する(時刻同期の関係上)
- projectのimport/export
- Client Secretに証明書を追加できるようにする
  - portalのアップデートだけでよい？
- ~~パスワード以外でのユーザーのログイン~~
  - ~~証明書~~
- ~~SAML対応~~

## Portal(Admin Console) enhancement

- 各ページの作成
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
  - customrole
    - update
  - audit event
  - project
    - reset secret
- outputの修正
  - error出力の方法
  - debug出力
- Production向け実行ファイルの作成

## new commands

- Gateway
  - Backendのユーザープログラムに対してアクセス制御するようなツール
  - keycloak-gatekeeperのようなものを想定

## operation enhancement

- add kubernetes yaml file
- add release pipeline
  - create binary files
- Benchmark
- exampleサイトの作成
