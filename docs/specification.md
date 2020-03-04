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
      - role
      - user
      - client
    - operate-type: 以下のいずれか
      - read
      - write
- custom role
  - Name: 3~63文字

## ユーザーパスワードロック

- x分以内にn回連続で間違ったパスワードを入力するとそのユーザーはy分ロックされる
- n回以内に成功すると失敗回数はリセットされる
- 変数値
  - n: LockCount
  - x: LockDuration
  - y: FailureResetTime
