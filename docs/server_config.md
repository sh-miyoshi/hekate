# サーバーの設定値

## 概要

このページでは、Hekateサーバーを起動時の設定値について解説します  
なお、有効となる優先順位はコマンドライン引数 > 環境変数 > configファイルとなっています。
つまり、コマンドライン引数で指定された値が最も優先して設定されます。

## 一覧

| 名称 | configファイル | 環境変数 | コマンドライン引数 | 説明 |
| :---- | :---- | :---- | :---- | :---- |
| configファイルパス | - | - | config | configファイルのパス |
| Adminユーザー名 | admin_name | HEKATE_ADMIN_NAME | admin | Adminユーザーの名前 |
| Adminユーザーパスワード | admin_password | HEKATE_ADMIN_PASSWORD | password | Adminユーザーのパスワード |
| サーバーポート | server_port | HEKATE_SERVER_PORT | port | サーバーを起動する際のポート番号 |
| サーバーバインドアドレス | server_bind_address | HEKATE_SERVER_BIND_ADDR | bind-addr | サーバーにバインドするアドレス |
| https有効化 | https.enabled | - | https | サーバーをhttpsで起動します |
| https証明書ファイルパス | https.cert-file | - | https-cert-file | httpsサーバー用の証明書ファイルのパス |
| https鍵ファイルパス | https.key-file | - | https-key-file | httpsサーバー用の鍵ファイルのパス |
| ログファイルパス | logfile | - | logfile | ログの出力先ファイルのパス。設定されてない、もしくは空文字列の場合は標準出力に表示されます |
| デバッグモード | debug_mode | HEKATE_ENV="DEBUG" | debug | デバッグ用のログも出力 |
| DBタイプ | db.type | HEKATE_DB_TYPE | db-type | サーバーが接続するDBのタイプ |
| DB接続文字列 | db.connection_string | HEKATE_DB_CONNECT_STRING | db-conn-str | DBに接続するための接続文字列 |
| ログインページリソースパス | user_login_page_res | HEKATE_LOGIN_PAGE_RES | login-res | ユーザーログインページのリソースへのパス |
