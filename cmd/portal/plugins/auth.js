import querystring from 'querystring'
import axios from 'axios'

export class AuthHandler {
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

  _setToken(obj) {
    const now = Date.now() / 1000
    window.localStorage.setItem('access_token', obj.access_token)
    window.localStorage.setItem('expires_in', now + obj.expires_in)
    window.localStorage.setItem('refresh_token', obj.refresh_token)
    window.localStorage.setItem(
      'refresh_expires_in',
      now + obj.refresh_expires_in
    )
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

  async _tokenRequest(opts) {
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
      return { ok: true, data: res.data }
    } catch (error) {
      console.log(error)
      if (error.response) {
        if (error.response.status >= 400 && error.response.status < 500) {
          return {
            ok: false,
            message: 'auth required',
            statusCode: error.response.status
          }
        }
      }
      return {
        ok: false,
        message: 'Failed to request the server',
        statusCode: 500
      }
    }
  }

  async AuthCode(authCode) {
    // TODO(consider state)
    const opts = {
      grant_type: 'authorization_code',
      client_id: 'admin-cli', // TODO(use param)
      code: authCode
    }

    const res = await this._tokenRequest(opts)
    if (res.ok) {
      console.log('successfully got token: %o', res.data)
      this._setToken(res.data)
      this.context.redirect('/home')
    } else if (res.statusCode >= 400 && res.statusCode < 500) {
      // redirect to login page
      this.context.redirect('/')
    } else {
      this.context.error({
        message: res.message,
        statusCode: 500
      })
    }
  }

  async TokenRefresh(refreshToken) {
    // TODO(consider state)
    const opts = {
      grant_type: 'refresh_token',
      client_id: 'admin-cli', // TODO(use param)
      refresh_token: refreshToken
    }

    const res = await this._tokenRequest(opts)
    if (res.ok) {
      console.log('successfully got token: %o', res.data)
      this._setToken(res.data)
    }

    return res
  }
}

export default (context, inject) => {
  inject('auth', new AuthHandler(context))
}
