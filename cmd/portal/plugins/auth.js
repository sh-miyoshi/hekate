import querystring from 'querystring'
import axios from 'axios'
import jwtdecode from 'jwt-decode'

export class AuthHandler {
  constructor(context) {
    this.context = context
    this.client_id = 'portal'
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

  _setLoginUser(token) {
    const data = jwtdecode(token)
    window.localStorage.setItem('user', data.preferred_username)
  }

  Login(project) {
    window.localStorage.setItem('login_project', project)

    // TODO(consider state)

    let protcol = 'https'
    if (!process.env.https) {
      protcol = 'http'
    }

    const opts = {
      scope: 'openid',
      response_type: 'code',
      client_id: this.client_id,
      redirect_uri:
        protcol +
        '://' +
        process.env.HEKATE_PORTAL_HOST +
        ':' +
        process.env.HEKATE_PORTAL_PORT +
        '/callback'
    }

    const url =
      process.env.HEKATE_SERVER_ADDR +
      '/api/v1/project/' +
      project +
      '/openid-connect/auth?' +
      this._encodeQuery(opts)

    window.location = url
  }

  async _tokenRequest(opts) {
    // TODO(use param: timeout)
    const headers = {
      'Content-Type': 'application/x-www-form-urlencoded'
    }

    const project = window.localStorage.getItem('login_project')
    const url =
      process.env.HEKATE_SERVER_ADDR +
      '/api/v1/project/' +
      project +
      '/openid-connect/token'

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
      client_id: this.client_id,
      code: authCode
    }

    const res = await this._tokenRequest(opts)
    if (res.ok) {
      console.log('successfully got token: %o', res.data)
      this._setLoginUser(res.data.access_token)
      this._setToken(res.data)
      this.context.redirect('/admin/home')
    } else if (res.statusCode >= 400 && res.statusCode < 500) {
      // redirect to login page
      this.context.redirect('/admin')
    } else {
      this.context.error({
        message: res.message,
        statusCode: 500
      })
    }
  }

  async GetToken() {
    const expiresIn = window.localStorage.getItem('expires_in')
    const refreshToken = window.localStorage.getItem('refresh_token')
    const now = Date.now() / 1000
    if (now < expiresIn) {
      return {
        ok: true,
        accessToken: window.localStorage.getItem('access_token')
      }
    }

    // TODO(consider state)
    const opts = {
      grant_type: 'refresh_token',
      client_id: this.client_id,
      refresh_token: refreshToken
    }

    const res = await this._tokenRequest(opts)
    if (res.ok) {
      console.log('successfully got token: %o', res.data)
      this._setToken(res.data)
      return {
        ok: true,
        accessToken: window.localStorage.getItem('access_token')
      }
    }

    return res
  }

  async Logout() {
    const refreshToken = window.localStorage.getItem('refresh_token')
    if (refreshToken) {
      const opts = {
        token_type_hint: 'refresh_token',
        refresh_token: refreshToken
      }

      const headers = {
        'Content-Type': 'application/x-www-form-urlencoded'
      }

      // TODO(use param: timeout)
      const url =
        process.env.HEKATE_SERVER_ADDR +
        '/api/v1/project/' +
        this.project +
        '/openid-connect/revoke'

      const handler = axios.create({
        headers,
        timeout: 10000
      })

      try {
        await handler.post(url, querystring.stringify(opts))
      } catch (error) {
        console.log(error)
      }
    }
    window.localStorage.removeItem('access_token')
    window.localStorage.removeItem('refresh_token')
    window.localStorage.removeItem('expires_in')
    window.localStorage.removeItem('refresh_expires_in')
    window.localStorage.removeItem('user')
    window.localStorage.removeItem('login_project')
  }
}

export default (context, inject) => {
  inject('auth', new AuthHandler(context))
}
