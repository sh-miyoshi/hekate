import querystring from 'querystring'
import axios from 'axios'

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

    // TODO(use param: project_name)
    const url =
      process.env.SERVER_ADDR +
      '/api/v1/project/master/openid-connect/auth?' +
      this._encodeQuery(opts)

    window.location = url
  }

  async AuthCode(authCode) {
    // TODO(consider state)
    const opts = {
      grant_type: 'authorization_code',
      client_id: 'admin-cli', // TODO(use param)
      code: authCode
    }

    // TODO(use param: project_name, timeout)
    const headers = {
      'Content-Type': 'application/x-www-form-urlencoded'
    }

    const url =
      process.env.SERVER_ADDR + '/api/v1/project/master/openid-connect/token'

    const handler = axios.create({
      headers,
      timeout: 10000
    })

    try {
      const res = await handler.post(url, querystring.stringify(opts))

      console.log(res)

      // set token to local storage and redirect to top page
    } catch (error) {
      console.log(error)
      if (error.response) {
        if (error.response.status >= 400 && error.response.status < 500) {
          // redirect to login page
          this.context.redirect('/')
          return
        }
      }
      this.context.error({
        message: 'Failed to request the server',
        statusCode: 500
      })
    }
  }
}

export default (context, inject) => {
  inject('auth', new AuthHandler(context))
}
