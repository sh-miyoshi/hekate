# TODO List

## application enhancement

- open id connect連携
- audit log
- projectのimport/export
- refresh tokenでの認証
- Gatewayプログラムの追加
  - keycloak-gatekeeperのようなものを想定
- GUI画面の追加
- 認可部分の関数化
  - httpパッケージとして追加
- 設定項目の追加
  - パスワードポリシー
  - refresh tokenのrevoke
  - cache
  - encrypt_type
- (project/user) enabledの有効化
- SAML対応
- テストの追加
  - ロジック部分のunit test
  - API部分のテスト

## operation enhancement

- add kubernetes yaml file
- write usage to README.md
- create public docker image
- configure CI
