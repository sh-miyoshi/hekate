# Dockerを使用してインストールする方法

このページでは[Docker](https://docs.docker.com/)を使用して`Hekate`をインストールする方法について述べます。
なお、あらかじめDockerの[インストールページ](https://docs.docker.com/install/)を参考にDockerをインストールしてください。

## サーバとポータルをAll in Oneで起動

```bash
# SERVER_ADDRはアクセスしたい場所からアクセスできるアドレスにしてください。
export SERVER_ADDR=localhost

# この値を指定していない場合はデフォルトの値(admin/password)が使用されます
export ADMIN_NAME=admin
export ADMIN_PASSWORD=password

# デフォルトのport番号以外をbindingする際は、以下の値もdocker起動時に環境変数で指定する必要があります
#  SERVER_PORT <- API_SERVER側のポート番号を変更したい場合
#  PORTAL_PORT <- PORTAL側のポート番号を変更したい場合
# 以下では例として設定していますが、デフォルトのポート番号(3000, 18443)を使用する場合は必要ありません
export PORTAL_PORT=3000
export SERVER_PORT=18443

docker run -d --name hekate \
  -p $PORTAL_PORT:3000 -p $SERVER_PORT:18443 \
  -e SERVER_ADDR=$SERVER_ADDR \
  -e HEKATE_ADMIN_NAME=$ADMIN_NAME \
  -e HEKATE_ADMIN_PASSWORD=$ADMIN_PASSWORD \
  -e PORTAL_PORT=$PORTAL_PORT \
  -e SERVER_PORT=$SERVER_PORT \
  smiyoshi/hekate:all-in-one
```

## ポータルへアクセス

[http://localhost:3000](http://localhost:3000)へアクセスし、先ほど指定した管理者名とパスワードでログインしてください。
