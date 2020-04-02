# APIレベルでのアクセス制御

このページではHekateを使用して、API単位であなたのサーバのアクセス制御をする方法を説明します。

## 全体像

TODO

## 前準備

README.mdを参考にHekateをインストールしてください。

以降では、Hekateには[http://localhost:3000](http://localhost:3000)でアクセスできるものとして説明します。
他の場合は適宜読み替えてください。

## 手順

### アクセスを制御したいサーバの準備

TODO

### Access Proxy用のOpenID Connect Clientの登録

- Portal([http://localhost:3000](http://localhost:3000))にアクセス
- Adminユーザの名前とパスワードを入力し、ログイン
- 左枠のメニューからClientを選択
- Add New Clinetボタンを押下
- Client IDを入力し、Createボタンを押下
  - ここではClient ID: `sample-gw`とする
- Client一覧画面から`sample-gw`のEditボタンを押下
- 表示されているSecretを記憶する

### Access Proxyの設置

今回はAccess Proxyとして[keycloak-gatekeeper](https://github.com/keycloak/keycloak-gatekeeper)を使用します。

- configファイルの準備

  ```bash
  # TODO
  ```

- keycloak-gatekeeperの起動
  - TODO
- アクセス
  - TODO
