# APIレベルでのアクセス制御

このページではHekateを使用して、API単位であなたのサーバのアクセス制御をする方法を説明します。

## 全体像

![イメージ図](../assets/access_ctrl_image.png)

## 前準備

### Hekateのインストール

[README.md](../../README.md)を参考にHekateをインストールしてください。

以降では、Hekateには以下のアドレスでアクセスできるものとして説明します。
他の場合は適宜読み替えてください。

- Portal: [http://localhost:3000](http://localhost:3000)
- Server: [http://localhost:18443](http://localhost:18443)

### アクセスを制御したいサーバ(リソースサーバ)の準備

アクセスを制御したい対象のサーバを起動します。
この例ではnginxを使用します。

```bash
docker run --name nginx -p 80:80 -d nginx
```

## 手順

### プロジェクトの準備

※新規プロジェクトを作成せずmasterプロジェクトでの操作も可能ですが、masterは管理者用プロジェクトなのでユーザー用の別プロジェクトを作成することをお勧めします

- Portal([http://localhost:3000](http://localhost:3000))にアクセス
- Adminユーザの名前とパスワードを入力し、ログイン
- 左枠のメニューからmasterを選択
- Add New Projectボタンを押下
- Nameを入力し、Createボタンを押下
  - ここでは、name: `sample`とする
- Project一覧から`sample`を選択
  - 左上が`master`から`sample`に変わっていることを確認

### ユーザーの追加

- 左枠のメニューからUserを選択
- Add New Userボタンを押下
- NameとPasswordを入力し、Createボタンを押下
  - ここではName: `admin`とする

### Access Proxy用のOpenID Connect Clientの登録

- 左枠のメニューからClientを選択
- Add New Clinetボタンを押下
- Client IDを入力し、Createボタンを押下
  - ここではClient ID: `sample-proxy`とする
- Client一覧画面から`sample-proxy`のEditボタンを押下
- 表示されているSecretを記憶する

### アクセス用のロールを作成・ユーザーに付与

- 左枠のメニューからRoleを選択
- Add New Roleボタンを押下
- Nameを入力し、Createボタンを押下
  - ここではName: `hello-access`とする
- 左枠のメニューからUserを選択
- ユーザ一覧からadminのEditボタンを押下
- Custom Roleで`hello-access`を選択しAssignボタンを押下
- Updateボタンを押下

### Access Proxyの設置

今回はAccess Proxyとして[keycloak-gatekeeper](https://github.com/keycloak/keycloak-gatekeeper)を使用します。

#### dockerコンテナを使用する場合

```bash
export CLIENT_SECRET="<確認したSecret>"
export HEKATE_SERVER="http://localhost:18443"
export RESOURCE_SERVER="http://localhost" # アクセスを制御したいサーバのアドレス

# configファイルの用意
cat << EOF > config.yaml
client-id: sample-proxy
client-secret: $CLIENT_SECRET
discovery-url: $HEKATE_SERVER/api/v1/project/sample # プロジェクトを変更した場合は適宜修正してください
enable-default-deny: true
skip-openid-provider-tls-verify: true
encryption_key: secret
listen: 0.0.0.0:5000
upstream-url: $RESOURCE_SERVER
secure-cookie: false
resources:
  - uri: /*
    methods:
    - GET
    roles:
      - user:hello-access
EOF

# keycloak-gatekeeperの起動
docker run --name gatekeeper -d --network host -v $PWD:/mnt/conf \
  quay.io/keycloak/keycloak-gatekeeper \
  --config=/mnt/conf/config.yaml
```

#### kuberentesを使用する場合

### アクセス

ブラウザから[http://localhost:5000](http://localhost:5000)にアクセス

TODO(dockerを使用している際にkeycloak-gatekeeperにアクセスできないとき)
TODO

### 後始末

TODO
