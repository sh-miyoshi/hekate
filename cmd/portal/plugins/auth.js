class AuthHandler {
  constructor(context) {
    this.context = context
  }

  _encodeQuery(queryObject) {
    return Object.entries(queryObject)
      .filter(([key, value]) => typeof value !== 'undefined')
      .map(
        ([key, value]) =>
          encodeURIComponent(key) +
          (value != null ? '=' + encodeURIComponent(value) : '')
      )
      .join('&')
  }

  Login() {
    // TODO(consider state)
    const opts = {
      scope: 'openid',
      response_type: 'code',
      client_id: 'admin-cli', // TODO(use param)
      redirect_uri: 'http://localhost:3000/callback' // TODO(use param)
    }

    const url =
      process.env.SERVER_ADDR +
      '/api/v1/project/master/openid-connect/auth?' +
      this._encodeQuery(opts)

    window.location = url
  }
}

export default (context, inject) => {
  inject('auth', new AuthHandler(context))
}
