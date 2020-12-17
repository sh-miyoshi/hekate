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

## (Optional)証明書の準備

ServerおよびPortalを独自の証明書で運用することができます。  
※container内に内蔵されているデフォルトのオレオレ証明書を使う場合や、前段でTLS終端するためhttpサーバーとして運用する場合はこの手順をスキップしてください。

```bash
# 証明書を作成
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=hekate.localhost"

# Kubernetes yamlの作成
crt_value=`base64 tls.crt | tr -d '\n'`
key_value=`base64 tls.key | tr -d '\n'`

cat << EOF > hekate-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: hekate-tls
type: kubernetes.io/tls
data:
  tls.crt: ${crt_value}
  tls.key: ${key_value}
EOF

# kubectl apply
kubectl apply -f hekate-secret.yaml
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
# 独自証明書を使用する場合はvolumesの部分とvolumeMountsの部分のコメントを外して有効化してください
kubectl apply -f server.yaml
```

## (Optional)Portalの構築

Portalのデプロイは必須ではありません。
もしCLIツールを使う場合など必要なければデプロイしなくても問題ありません。
また、Portalはサーバー管理者のみ利用できればいいのでアクセス制限等を実施しても構いません。

```bash
wget https://github.com/sh-miyoshi/hekate/raw/master/deployments/kubernetes/portal.yaml
vi portal.yaml
# https://localhost:3000以外のアドレスを使用する場合は適宜修正してください。
# 独自証明書を使用する場合はvolumesの部分とvolumeMountsの部分のコメントを外して有効化してください

kubectl apply -f portal.yaml
```

Pod起動後、[https://localhost:3000](https://localhost:3000)でアクセスできます。
ただし、portalの起動にはビルド処理を含むため少し時間がかかります。
`kubectl logs`コマンドなどでログを確認しつつ、気長に待ってください。

TODO(現在執筆中です・・・)
