# コマンドラインツールの仕様

## 概要

- 名前: jwtctl

## 引数に関して

基本的に1wordで対応
どうしても複数wordになる場合は`-`で区切る
配列を引数にする場合はStringSliceメソッドを使用する

### グローバル引数

- `--output`, `-o`: 出力フォーマットを変更する(text, json)
- `--debug`: (WIP)debugログを出力する

## ログ・メッセージの仕様

- メッセージタイプ
  - Debug
  - Info
  - Error
- Info, Debugメッセージは標準出力に、Errorメッセージは標準エラー出力に表示する
- Errorメッセージの場合、表示する際は`[ERROR]`とつける
