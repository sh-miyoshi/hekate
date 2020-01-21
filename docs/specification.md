# 仕様書

## 各リソースの制限事項

- project
  - Name: 2~31文字 && (英語小文字 || 数字 || -) && 先頭文字は英語小文字
  - TokenConfig
    - AccessTokenLifeSpan: 1 sec以上
    - RefreshTokenLifeSpan: 1 sec以上
    - SigningAlgorithm: 以下のいずれかであること
      - RS256
- client
  - ID: 2~127文字
  - ProjectName: project.Nameと同じ
  - Secret: 8~255文字
  - AccessType: public || confidential
  - AllowedCallbackURLs: URL形式の配列であること
- user
  - ProjectName: project.Nameと同じ
  - Name: 3~63文字
  - Roles: role形式(\<resource-type\>-\<operate-type\>)の配列であること
    - resource-type: 以下のいずれか
      - cluster
      - project
      - role
      - user
      - client
    - operate-type: 以下のいずれか
      - read
      - write
