# APIレベルでのアクセス制御

このページではHekateを使用して、API単位であなたのサーバのアクセス制御をする方法を説明します。

## 全体像

TODO

## 前準備

README.mdを参考にHekateをインストールしてください。

以降では、Hekateには以下のアドレスでアクセスできるものとして説明します。
他の場合は適宜読み替えてください。

- Portal: [http://localhost:3000](http://localhost:3000)
- Server: [http://localhost:8080](http://localhost:8080)

## 手順

### アクセスを制御したいサーバの準備

アクセスを制御したい対象のサーバを起動します。
この例では[test-server](https://github.com/sh-miyoshi/test-server)を使用します。

```bash
docker run --name test-server -p 10000:10000 -d smiyoshi/test-server
```

### Access Proxy用のOpenID Connect Clientの登録

- Portal([http://localhost:3000](http://localhost:3000))にアクセス
- Adminユーザの名前とパスワードを入力し、ログイン
- 左枠のメニューからClientを選択
- Add New Clinetボタンを押下
- Client IDを入力し、Createボタンを押下
  - ここではClient ID: `sample-gw`とする
- Client一覧画面から`sample-gw`のEditボタンを押下
- 表示されているSecretを記憶する

### アクセス用のロールを作成・ユーザーに付与

- Potalにログイン後、左枠のメニューからRoleを選択
- Add New Roleボタンを押下
- Nameを入力し、Createボタンを押下
  - ここではName: `hello-access`とする
- 左枠のメニューからUserを選択
- ユーザ一覧からadminのEditボタンを押下
- Custom Roleで`hello-access`を選択しAssignボタンを押下
- Updateボタンを押下

### Access Proxyの設置

今回はAccess Proxyとして[keycloak-gatekeeper](https://github.com/keycloak/keycloak-gatekeeper)を使用します。

```bash
# configファイルの用意
cat << EOF > config.yaml
client-id: sample-gw
client-secret: <確認したSecret>
discovery-url: http://localhost:8080/api/v1/project/master # Hekateサーバのアドレスとプロジェクトを変更した場合は適宜修正してください
enable-default-deny: true
skip-openid-provider-tls-verify: true
encryption_key: secret
listen: 0.0.0.0:5000
upstream-url: http://localhost:10000 # アクセスを制御したいサーバのアドレス
secure-cookie: false
resources:
  - uri: /hello
    methods:
    - GET
    roles:
      - user:hello-access
EOF

# keycloak-gatekeeperの起動
# docker run -it --rm quay.io/keycloak/keycloak-gatekeeper

# アクセス
```
