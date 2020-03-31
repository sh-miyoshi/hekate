# Dockerを使用してインストールする方法

このページでは[Docker](https://docs.docker.com/)を使用して`Hekate`をインストールする方法について述べます。
なお、あらかじめDockerの[インストールページ](https://docs.docker.com/install/)を参考にDockerをインストールしてください。

## サーバとポータルをAll in Oneで起動

※現在設定ファイル的な理由で3000, 8080番以外のポートをport bindingできません。

```bash
# SERVER_ADDRはアクセスしたい場所からアクセスできるアドレスにしてください。
export SERVER_ADDR=localhost

# この値を指定していない場合はデフォルトの値(admin/password)が使用されます
export ADMIN_NAME=admin
export ADMIN_PASSWORD=password

docker run -d --name hekate \
  -p 3000:3000 -p 8080:8080 \
  -e SERVER_ADDR=$SERVER_ADDR \
  -e HEKATE_ADMIN_NAME=$ADMIN_NAME \
  -e HEKATE_ADMIN_PASSWORD=$ADMIN_PASSWORD \
  smiyoshi/hekate:all-in-one
```

## ポータルへアクセス

[http://localhost:3000](http://localhost:3000)へアクセスし、先ほど指定した管理者名とパスワードでログインしてください。
