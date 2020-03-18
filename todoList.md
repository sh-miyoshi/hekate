# TODO List

## Bugs(?)

- userが削除された場合のaccess tokenが無効化されない(この仕様はok?)
- custom role削除時にユーザからremoveしなくてよい？
- 別projectのリソースを変更できないか要確認

## new commands

- Gateway
  - Backendのユーザープログラムに対してアクセス制御するようなツール
  - keycloak-gatekeeperのようなものを想定

## server application enhancement

- userのパスワード変更のrole見直し
  - 本人のみが変更できるようにする
- API responseのエラーコードが足りてないバグの修正
  - ClientUpdateHandler
  - RoleUpdateHandler
  - UserUpdateHandler
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
- DBGCの追加
  - Expiredしたsessionなどを一定期間ごとに削除する
- リソース削除時に紐づくリソースも削除
- db manager validationの追加
- filterの追加(user, role)
- audit log
  - time
  - resource type (or url path and method)
  - client
  - success or failed
- 各種APIの実装
  - sessionの詳細取得(引数: project, useID, sessionID)
- テストの追加
  - ロジック部分のunit test
  - API部分のテスト
    - エラー処理回り
- 設定項目の追加
  - パスワードポリシー
  - encrypt_type(signing_method)
- http errorの充実
  - example: [facebook for developers](https://developers.facebook.com/docs/messenger-platform/reference/send-api/error-codes?locale=ja_JP)
- projectのimport/export
- OpenID Connect部分のエンハンス
  - openid connect APIの実装
    - implicit flow
    - hybrid flow
  - TokenAPIでredirect_uriのチェック
    - query内に存在するならallowed listにあるかチェック
  - access tokenのrevocation
  - AuthRequestに他のパラメータを追加
  - code認証失敗時、すべてのtokenを無効化
  - subject_types_supportedにpairwiseをサポート
  - RS256以外のSigining Algorithmのサポート
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
- API responseのtime formatの見直し

## Portal(Admin Console) enhancement

- middleware処理
  - roleが足りない(masterプロジェクトにいない、cluster操作権限がない?)
    - 強制ログアウト or 白紙のページを見せる(こっちが有力)
- redirect先
  - 前回開いていたページに戻る
- alert画面のcss修正
- project選択を一覧に
- logoutボタンの実装
- headerにユーザー名を表示
- 各ページの作成
  - User
    - user削除コマンドの実装
    - user createページの作成
    - user更新コマンドの実装
  - Role
  - Client
- home pageにmiddlewareを使用
- oidc authのstateチェック
- client_idやproject nameをどうするか(変数化)
- 各ユーザーのアカウント管理画面
  - user名変更
  - パスワード変更

## CLI tool(hctl) enhancement

- configファイルの扱い
  - `no such file or directory`のとき、新規作成
  - permission: 0700, 0600
  - デフォルト値(localhost:8080, master)
- 各APIへの対応
  - project
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

## operation enhancement

- add kubernetes yaml file
- write usage to README.md
- add release pipeline
  - create public docker image
  - create binary files

## For production

- 初期allowed callback urlの初期化
- CLIのinsecure, timeoutの設定見直し
