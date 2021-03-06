import axios from 'axios'
import { AuthHandler } from './auth'

class APIClient {
  constructor(context) {
    this.handler = new AuthHandler(context)
    this.serverAddr = process.env.HEKATE_SERVER_ADDR
  }

  async _request(url, method, data) {
    let res = await this.handler.GetToken()

    if (!res.ok) {
      return res
    }

    const headers = {
      Authorization: 'Bearer ' + res.accessToken
    }

    if (!data) {
      headers['Content-Type'] = 'application/json'
    }

    const handler = axios.create({
      headers,
      timeout: 10000
    })

    try {
      switch (method) {
        case 'GET':
          res = await handler.get(url)
          break
        case 'POST':
          res = await handler.post(url, data)
          break
        case 'PUT':
          res = await handler.put(url, data)
          break
        case 'DELETE':
          res = await handler.delete(url)
          break
        default:
          return {
            ok: false,
            message: 'HTTP Method ' + method + ' is unsupported',
            statusCode: 500
          }
      }
      return { ok: true, data: res.data }
    } catch (error) {
      console.log(error)
      if (error.response) {
        if (error.response.status >= 400 && error.response.status < 500) {
          let msg = error.response.data.error
          if (
            error.response.data.error_description != null &&
            error.response.data.error_description.length > 0
          ) {
            msg = error.response.data.error_description
          }

          return {
            ok: false,
            message: msg,
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

  async ProjectGetList() {
    const url = this.serverAddr + '/adminapi/v1/project'
    const res = await this._request(url, 'GET')
    return res
  }

  async ProjectGet(projectName) {
    const url = this.serverAddr + '/adminapi/v1/project/' + projectName
    const res = await this._request(url, 'GET')
    return res
  }

  async ProjectCreate(projectName) {
    const data = {
      name: projectName,
      token_config: {
        access_token_life_span: 300, // 5 minutes
        refresh_token_life_span: 1209600, // 2 weeks
        signing_algorithm: 'RS256'
      },
      allow_grant_types: [
        // set default allow grant types
        'authorization_code',
        'client_credentials',
        'refresh_token'
      ]
    }
    const url = this.serverAddr + '/adminapi/v1/project'
    const res = await this._request(url, 'POST', data)
    return res
  }

  async ProjectDelete(projectName) {
    const url = this.serverAddr + '/adminapi/v1/project/' + projectName
    const res = await this._request(url, 'DELETE')
    return res
  }

  async ProjectUpdate(projectName, info) {
    const url = this.serverAddr + '/adminapi/v1/project/' + projectName
    const res = await this._request(url, 'PUT', info)
    return res
  }

  async UserCreate(projectName, info) {
    const url =
      this.serverAddr + '/adminapi/v1/project/' + projectName + '/user'
    const res = await this._request(url, 'POST', info)
    return res
  }

  async UserGetList(projectName) {
    const url =
      this.serverAddr + '/adminapi/v1/project/' + projectName + '/user'
    const res = await this._request(url, 'GET')
    return res
  }

  async UserGet(projectName, userID) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/user/' +
      userID
    const res = await this._request(url, 'GET')
    return res
  }

  async UserDelete(projectName, userID) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/user/' +
      userID
    const res = await this._request(url, 'DELETE')
    return res
  }

  async UserUpdate(projectName, userID, info) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/user/' +
      userID
    const res = await this._request(url, 'PUT', info)
    return res
  }

  async UserUnlock(projectName, userID) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/user/' +
      userID +
      '/unlock'
    const res = await this._request(url, 'POST')
    return res
  }

  async UserResetPassword(projectName, userID, newPassword) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/user/' +
      userID +
      '/reset-password'
    const res = await this._request(url, 'POST', { password: newPassword })
    return res
  }

  async ClientCreate(projectName, info) {
    const url =
      this.serverAddr + '/adminapi/v1/project/' + projectName + '/client'
    const res = await this._request(url, 'POST', info)
    return res
  }

  async ClientGetList(projectName) {
    const url =
      this.serverAddr + '/adminapi/v1/project/' + projectName + '/client'
    const res = await this._request(url, 'GET')
    return res
  }

  async ClientGet(projectName, clientID) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/client/' +
      clientID
    const res = await this._request(url, 'GET')
    return res
  }

  async ClientDelete(projectName, clientID) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/client/' +
      clientID
    const res = await this._request(url, 'DELETE')
    return res
  }

  async ClientUpdate(projectName, clientID, info) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/client/' +
      clientID
    const res = await this._request(url, 'PUT', info)
    return res
  }

  async RoleCreate(projectName, info) {
    const url =
      this.serverAddr + '/adminapi/v1/project/' + projectName + '/role'
    const res = await this._request(url, 'POST', info)
    return res
  }

  async RoleGetList(projectName) {
    const url =
      this.serverAddr + '/adminapi/v1/project/' + projectName + '/role'
    const res = await this._request(url, 'GET')
    return res
  }

  async RoleGet(projectName, roleID) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/role/' +
      roleID
    const res = await this._request(url, 'GET')
    return res
  }

  async RoleDelete(projectName, roleID) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/role/' +
      roleID
    const res = await this._request(url, 'DELETE')
    return res
  }

  async RoleUpdate(projectName, roleID, info) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/role/' +
      roleID
    const res = await this._request(url, 'PUT', info)
    return res
  }

  async SessionGet(projectName, sessionID) {
    const url =
      this.serverAddr +
      '/adminapi/v1/project/' +
      projectName +
      '/session/' +
      sessionID
    const res = await this._request(url, 'GET')
    return res
  }

  async KeysGet(projectName) {
    const url =
      this.serverAddr + '/adminapi/v1/project/' + projectName + '/keys'
    const res = await this._request(url, 'GET')
    return res
  }

  async KeysReset(projectName) {
    const url =
      this.serverAddr + '/adminapi/v1/project/' + projectName + '/keys/reset'
    const res = await this._request(url, 'POST')
    return res
  }

  async AuditGetList(projectName) {
    // TODO(from, to, pagenation)
    const url =
      this.serverAddr + '/adminapi/v1/project/' + projectName + '/audit'
    const res = await this._request(url, 'GET')
    return res
  }

  // ----------------------------
  // User API Client
  // ----------------------------

  async UserAPIGetUser(projectName, userID) {
    const url =
      this.serverAddr + '/userapi/v1/project/' + projectName + '/user/' + userID
    const res = await this._request(url, 'GET')
    return res
  }

  async UserAPIChangePassword(projectName, userID, newPassword) {
    const url =
      this.serverAddr +
      '/userapi/v1/project/' +
      projectName +
      '/user/' +
      userID +
      '/change-password'
    const res = await this._request(url, 'POST', { password: newPassword })
    return res
  }

  async UserAPIOTPGenerate(projectName, userID) {
    const url =
      this.serverAddr +
      '/userapi/v1/project/' +
      projectName +
      '/user/' +
      userID +
      '/otp'
    const res = await this._request(url, 'POST')
    return res
  }

  async UserAPIOTPVerify(projectName, userID, userCode) {
    const url =
      this.serverAddr +
      '/userapi/v1/project/' +
      projectName +
      '/user/' +
      userID +
      '/otp/verify'
    const res = await this._request(url, 'POST', { user_code: userCode })
    return res
  }

  async UserAPIOTPDelete(projectName, userID) {
    const url =
      this.serverAddr +
      '/userapi/v1/project/' +
      projectName +
      '/user/' +
      userID +
      '/otp'
    const res = await this._request(url, 'DELETE')
    return res
  }
}

export default (context, inject) => {
  inject('api', new APIClient(context))
}
