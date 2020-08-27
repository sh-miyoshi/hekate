# All in One Image

All in Oneイメージではserverとportalを1つのdockerイメージで起動できます。
また初期DBとしてはメモリを使用するのでほかにDBを立てる必要がなく、テスト環境として使用しやすい形になっています。

## 内部のディレクトリ構成

```text
/hekate
|-- portal
|   `-- *) portal codes
|-- run.sh
|-- secret
|   |-- tls.crt
|   `-- tls.key
`-- server
    |-- config.yaml
    |-- hekate-server
    `-- *) login page assets
```
