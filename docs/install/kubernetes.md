# Kubernetesを使用してインストールする方法

このページでは[Kubernetes](https://kubernetes.io/ja/)を使用して`Hekate`をインストールする方法について述べます。  
このページの内容を実行する前にあらかじめKubernetesクラスターを構築してください。  
※Kubernetesクラスターを構築するもっとも簡単な方法は各クラウドプロバイダーのマネージドサービス(EKS, GKE, AKSなど)を使用することです。

## DBの構築

今回はMongo DBにデータを保存します。  
※Production環境の場合はマネージドサービスを使用することを強く推奨します。

```bash
kubectl apply -f https://github.com/sh-miyoshi/hekate/raw/master/deployments/kubernetes/mongo-pv.yaml
kubectl apply -f https://github.com/sh-miyoshi/hekate/raw/master/deployments/kubernetes/mongo.yaml
```

## Serverの構築

Mongo DBへの接続文字列やサーバー管理者のパスワードなど必要に応じてyamlファイルを修正してください。  
hekate-portalも立てる場合は先にportalにアクセスさせるためのアドレスを決めておき、server起動時に環境変数(HEKATE_PORTAL_ADDR)として渡す必要があります。

> hekateではportalをOpenID ConnectのClientとして実装しているため、serverに正しいClientのredirect_urlを登録する必要があります。この環境変数を指定することで管理プロジェクト(master)にportalのアドレスをredirect_urlとして登録します。

```bash
wget https://github.com/sh-miyoshi/hekate/raw/master/deployments/kubernetes/server.yaml
vi server.yaml
# db_conn_str, admin_passwordなどの秘密情報を修正
# portalもデプロイする場合はHEKATE_PORTAL_ADDRも設定
kubectl apply -f server.yaml
```

## Portalの構築

Portalのデプロイは必須ではありません。
もしCLIツールを使う場合など必要なければデプロイしなくても問題ありません。
また、Portalはサーバー管理者のみ利用できればいいのでアクセス制限等を実施しても構いません。

```bash
wget https://github.com/sh-miyoshi/hekate/raw/master/deployments/kubernetes/portal.yaml
```

TODO(現在執筆中です・・・)
