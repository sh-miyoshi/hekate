import querystring from 'querystring'
import crypto from 'crypto'
import axios from 'axios'
import jwtdecode from 'jwt-decode'
import base64url from 'base64url'

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

  _setLoginUser(token) {
    const data = jwtdecode(token)
    window.localStorage.setItem('user_name', data.preferred_username)
    window.localStorage.setItem('user_id', data.sub)
  }

  _genRandom(length) {
    // 1 byte = 2 characters in hex
    const buf = crypto.randomBytes(length / 2)
    return buf.toString('hex')
  }

  _createCodeChallenge(verifier, method) {
    switch (method) {
      case 'PLANE':
        return verifier
      case 'S256': {
        const hash = crypto
          .createHash('sha256')
          .update(verifier, 'ascii')
          .digest()
        return base64url.encode(hash)
      }
    }
    console.log('Invalid code challnge method ', method)
    return ''
  }

  _setRedirectTo() {
    let url = '/admin'
    const path = this.context.from.path
    if (path.startsWith('/user/')) {
      const re = /\/user\/project\/[^/]+/g
      const found = path.match(re)
      if (found.length > 0) {
        url = found[0] + '/info'
      }
    }

    window.localStorage.setItem('redirect_to', url)
  }

  RemoveAllData() {
    // remove cookie data
    document.cookie = 'HEKATE_LOGIN_SESSION=; max-age=0'

    // remove local storage data
    window.localStorage.removeItem('access_token')
    window.localStorage.removeItem('refresh_token')
    window.localStorage.removeItem('expires_in')
    window.localStorage.removeItem('refresh_expires_in')
    window.localStorage.removeItem('user_name')
    window.localStorage.removeItem('user_id')
    window.localStorage.removeItem('login_project')
    window.localStorage.removeItem('redirect_to')
  }

  Login(project) {
    this._setRedirectTo()

    window.localStorage.setItem('login_project', project)

    const state = this._genRandom(8)
    window.sessionStorage.setItem('login_state', state)

    const verifier = this._genRandom(128)
    window.sessionStorage.setItem('code_verifier', verifier)
    const challenge = this._createCodeChallenge(verifier, 'S256')

    const opts = {
      scope: 'openid email',
      response_type: 'code',
      client_id: process.env.CLIENT_ID,
      redirect_uri: process.env.HEKATE_PORTAL_ADDR + '/callback',
      code_challenge: challenge,
      code_challenge_method: 'S256',
      state
    }

    const url =
      process.env.HEKATE_SERVER_ADDR +
      '/authapi/v1/project/' +
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
      '/authapi/v1/project/' +
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

  async AuthCode(authCode, state) {
    const redirect = window.localStorage.getItem('redirect_to')

    const correctState = window.sessionStorage.getItem('login_state')
    console.log('state: ', correctState, ', received state: ', state)
    if (correctState !== state) {
      this.context.error({
        message:
          'The state parameters are different. It may have been attacked by CSRF.',
        statusCode: 500
      })
      return
    }
    const verifier = window.sessionStorage.getItem('code_verifier')

    const opts = {
      grant_type: 'authorization_code',
      client_id: process.env.CLIENT_ID,
      code: authCode,
      code_verifier: verifier,
      state
    }

    const res = await this._tokenRequest(opts)
    if (res.ok) {
      console.log('successfully got token: %o', res.data)
      this._setLoginUser(res.data.access_token)
      this._setToken(res.data)
      this.context.redirect(redirect)
    } else if (res.statusCode >= 400 && res.statusCode < 500) {
      // redirect to login page
      this.context.redirect(redirect)
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
      client_id: process.env.CLIENT_ID,
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
      const headers = {
        'Content-Type': 'application/x-www-form-urlencoded',
        Authorization: 'Bearer ' + window.localStorage.getItem('access_token')
      }

      // TODO(use param: timeout)
      const userID = window.localStorage.getItem('user_id')
      const project = window.localStorage.getItem('login_project')
      const url =
        process.env.HEKATE_SERVER_ADDR +
        '/userapi/v1/project/' +
        project +
        '/user/' +
        userID +
        '/logout'

      const handler = axios.create({
        headers,
        timeout: 10000
      })

      try {
        await handler.post(url)
      } catch (error) {
        console.log(error)
      }
    }
  }

  GetUserSystemRoles() {
    const token = window.localStorage.getItem('access_token')
    if (!token) {
      console.log('Failed to get access_token from local storage')
      return []
    }

    const user = jwtdecode(token)
    const res = user.resource_access.system_management.roles
    if (!res) {
      return []
    }
    return res
  }
}

export default (context, inject) => {
  inject('auth', new AuthHandler(context))
}
