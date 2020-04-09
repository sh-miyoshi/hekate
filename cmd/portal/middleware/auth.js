import { AuthHandler } from '../plugins/auth.js'

export default async function(context) {
  let redirect = '/'
  const project = window.localStorage.getItem('login_project')
  if (project && project !== 'master') {
    redirect = '/user/project/' + project
  }

  const expiresIn = window.localStorage.getItem('expires_in')
  if (!expiresIn) {
    context.redirect(redirect)
    return
  }

  const handler = new AuthHandler(context)
  const res = await handler.GetToken()
  if (!res.ok) {
    if (res.statusCode >= 400 && res.statusCode < 500) {
      context.redirect(redirect)
    } else {
      context.error({
        message: res.message,
        statusCode: 500
      })
    }
  }
}
