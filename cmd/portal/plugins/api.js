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
          return {
            ok: false,
            message: error.response.data.error,
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
    const url = this.serverAddr + '/api/v1/project'
    const res = await this._request(url, 'GET')
    return res
  }

  async ProjectGet(projectName) {
    const url = this.serverAddr + '/api/v1/project/' + projectName
    const res = await this._request(url, 'GET')
    return res
  }

  async ProjectCreate(projectName) {
    const data = {
      name: projectName,
      tokenConfig: {
        accessTokenLifeSpan: 300, // 5 minutes
        refreshTokenLifeSpan: 1209600, // 2 weeks
        signingAlgorithm: 'RS256'
      },
      allowGrantTypes: [
        // set default allow grant types
        'authorization_code',
        'client_credentials',
        'refresh_token'
      ]
    }
    const url = this.serverAddr + '/api/v1/project'
    const res = await this._request(url, 'POST', data)
    return res
  }

  async ProjectDelete(projectName) {
    const url = this.serverAddr + '/api/v1/project/' + projectName
    const res = await this._request(url, 'DELETE')
    return res
  }

  async ProjectUpdate(projectName, info) {
    const url = this.serverAddr + '/api/v1/project/' + projectName
    const res = await this._request(url, 'PUT', info)
    return res
  }

  async UserCreate(projectName, info) {
    const url = this.serverAddr + '/api/v1/project/' + projectName + '/user'
    const res = await this._request(url, 'POST', info)
    return res
  }

  async UserGetList(projectName) {
    const url = this.serverAddr + '/api/v1/project/' + projectName + '/user'
    const res = await this._request(url, 'GET')
    return res
  }

  async UserGet(projectName, userID) {
    const url =
      this.serverAddr + '/api/v1/project/' + projectName + '/user/' + userID
    const res = await this._request(url, 'GET')
    return res
  }

  async UserDelete(projectName, userID) {
    const url =
      this.serverAddr + '/api/v1/project/' + projectName + '/user/' + userID
    const res = await this._request(url, 'DELETE')
    return res
  }

  async UserUpdate(projectName, userID, info) {
    const url =
      this.serverAddr + '/api/v1/project/' + projectName + '/user/' + userID
    const res = await this._request(url, 'PUT', info)
    return res
  }

  async UserUnlock(projectName, userID) {
    const url =
      this.serverAddr +
      '/api/v1/project/' +
      projectName +
      '/user/' +
      userID +
      '/unlock'
    const res = await this._request(url, 'POST')
    return res
  }

  async ClientCreate(projectName, info) {
    const url = this.serverAddr + '/api/v1/project/' + projectName + '/client'
    const res = await this._request(url, 'POST', info)
    return res
  }

  async ClientGetList(projectName) {
    const url = this.serverAddr + '/api/v1/project/' + projectName + '/client'
    const res = await this._request(url, 'GET')
    return res
  }

  async ClientGet(projectName, clientID) {
    const url =
      this.serverAddr + '/api/v1/project/' + projectName + '/client/' + clientID
    const res = await this._request(url, 'GET')
    return res
  }

  async ClientDelete(projectName, clientID) {
    const url =
      this.serverAddr + '/api/v1/project/' + projectName + '/client/' + clientID
    const res = await this._request(url, 'DELETE')
    return res
  }

  async ClientUpdate(projectName, clientID, info) {
    const url =
      this.serverAddr + '/api/v1/project/' + projectName + '/client/' + clientID
    const res = await this._request(url, 'PUT', info)
    return res
  }

  async RoleCreate(projectName, info) {
    const url = this.serverAddr + '/api/v1/project/' + projectName + '/role'
    const res = await this._request(url, 'POST', info)
    return res
  }

  async RoleGetList(projectName) {
    const url = this.serverAddr + '/api/v1/project/' + projectName + '/role'
    const res = await this._request(url, 'GET')
    return res
  }

  async RoleGet(projectName, roleID) {
    const url =
      this.serverAddr + '/api/v1/project/' + projectName + '/role/' + roleID
    const res = await this._request(url, 'GET')
    return res
  }

  async RoleDelete(projectName, roleID) {
    const url =
      this.serverAddr + '/api/v1/project/' + projectName + '/role/' + roleID
    const res = await this._request(url, 'DELETE')
    return res
  }

  async RoleUpdate(projectName, roleID, info) {
    const url =
      this.serverAddr + '/api/v1/project/' + projectName + '/role/' + roleID
    const res = await this._request(url, 'PUT', info)
    return res
  }

  async SessionGet(projectName, sessionID) {
    const url =
      this.serverAddr +
      '/api/v1/project/' +
      projectName +
      '/session/' +
      sessionID
    const res = await this._request(url, 'GET')
    return res
  }

  async KeysGet(projectName) {
    const url = this.serverAddr + '/api/v1/project/' + projectName + '/keys'
    const res = await this._request(url, 'GET')
    return res
  }

  async KeysReset(projectName) {
    const url =
      this.serverAddr + '/api/v1/project/' + projectName + '/keys/reset'
    const res = await this._request(url, 'POST')
    return res
  }

  async AuditGetList(projectName) {
    // TODO(from, to, pagenation)
    const url = this.serverAddr + '/api/v1/project/' + projectName + '/audit'
    const res = await this._request(url, 'GET')
    return res
  }
}

export default (context, inject) => {
  inject('api', new APIClient(context))
}
