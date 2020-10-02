# 仕様書

## 各リソースの制限事項

- project
  - Name: 3~63文字 && (英語小文字 or 数字 or -._ ) && 先頭文字は英語小文字
  - TokenConfig
    - AccessTokenLifeSpan: 1[sec]以上
    - RefreshTokenLifeSpan: 1[sec]以上
    - SigningAlgorithm: 以下のいずれかであること
      - RS256
- client
  - ID: 3~63文字 && (英語小文字 or 数字 or -._ ) && 先頭文字は英語小文字
  - ProjectName: project.Nameと同じ
  - Secret: 8~255文字
  - AccessType: public || confidential
  - AllowedCallbackURLs: URL形式の配列であること
- user
  - ProjectName: project.Nameと同じ
  - Name: 3~63文字
  - System Roles: role形式(\<resource-type\>-\<operate-type\>)の配列であること
    - resource-type: 以下のいずれか
      - cluster
      - project
    - operate-type: 以下のいずれか
      - read
      - write
- custom role
  - Name: 3~63文字

## System Role

System Roleはcluter roleとproject roleがあり、それぞれに対してread, writeが存在する。

cluster roleはCluster全体を管理でき、master projectに存在するユーザーにのみ付与できる。
cluster権限はすべてのprojectのすべてのリソースに対してread, write可能である。

project roleは自プロジェクトのみに対して有効であり、他プロジェクトのリソースに干渉できない。
また、userに対してwriteのみを付与することはできない。
そのため、以下の制約をつける。

- Role 追加時
  - write権限を付与する際はすでにuserにそのリソースに対するread権限が存在するか、同時に登録しようとしているRoleの中にそのリソースに対するread権限が必要
- Role 削除時
  - read権限を削除する場合、userにそのリソースに対するwrite権限がついていないか、そのリソースのread権限も同時に削除しなければならない

## ユーザーパスワードロック

- x分以内にn回連続で間違ったパスワードを入力するとそのユーザーはy分ロックされる
- n回以内に成功すると失敗回数はリセットされる
- 変数値
  - n: LockCount
  - x: LockDuration
  - y: FailureResetTime
