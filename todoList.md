# TODO List

## application enhancement

- set config from text(yaml) file
  - e.g. `jwt-server --config=config.yaml`
- run with custom jwt token format
  - idea1: write jwt_return.json and jwt-server read it dynamically
    - merit: user can create any form of json
    - demerit: how to decide necessary field, and how to deal with it in golang
  - idea2: write token_config.yaml and generate token_return.json from yaml
    - merit: easy to control in golang
    - demerit: user cannot create custom field
- add user and role
- multi tenant(?)
- integrate to open id connect
- audit log
- event hook
- import/export setting(?)
- refresh token

## keycloak

- project
  - settings
    - name, enabled, endpoints(open id connect, saml, ...)
    - encrypt_type
    - cache
    - token
      - timeout, refresh_token, offline_token, revoke
  - clients(users?)
  - roles
  - user federation
  - authentication
    - password policy
  - audit events
    - config
  - import/export


## operation enhancement

- add kubernetes yaml file
- write usage to README.md
- create public docker image
- configure CI
