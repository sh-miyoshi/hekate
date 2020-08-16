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

## Portalの構築

TODO(現在執筆中です・・・)
