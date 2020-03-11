# TODO List

## Bugs(?)

- userが削除された場合のaccess tokenが無効化されない(この仕様はok?)
- custom role削除時にユーザからremoveしなくてよい？

## new commands

- Gateway
  - Backendのユーザープログラムに対してアクセス制御するようなツール
  - keycloak-gatekeeperのようなものを想定

## server application enhancement

- user auth failed処理のアップデート
  - loginSession structにScope,ResponseTypeを追加
  - user.LoginVerifyする前にsession情報を取得
  - 失敗時はsession情報から再度codeを発行
- PUT APIの修正
  - フィールドがない or nullの場合は更新しない
- Auhtorization Code Flowのエラー時処理の修正
  - error responseをhtmlで返す
- GET LIST APIの修正
  - 残り: client, role
  - すべての検索結果を取得するようにする
    - http queryでfilterできるようにする
      - API handlerの修正
      - DB queryの修正
  - API docの修正
    - 残り: client, role, project, user
- ユーザーパスワードロック
  - Project Infoに追加
  - APIの修正
    - ログインリクエスト失敗時にロックする
    - 強制ロック解除用のAPI
  - API docの修正
- DBGCの追加
  - Expiredしたsessionなどを一定期間ごとに削除する
- audit log
  - time
  - resource type (or url path and method)
  - client
  - success or failed
- 各種APIの実装
  - openid connect API
    - implicit flow
    - hybrid flow
  - sessionの詳細取得(引数: project, useID, sessionID)
- テストの追加
  - ロジック部分のunit test
  - API部分のテスト
    - エラー処理回り
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

- CSSの変更
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
  - Role
  - Client

## CLI tool(hctl) enhancement

- 各APIへの対応
  - project
    - update
  - user
    - update
    - role add
    - role delete
    - session revoke
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
- outputの修正
  - error出力の方法
  - debug出力
  - error messageの内容
- default config pathの修正
- configコマンドの作成・修正
- Production向け実行ファイルの作成

## operation enhancement

- add kubernetes yaml file
- write usage to README.md
- add release pipeline
  - create public docker image
  - create binary files
