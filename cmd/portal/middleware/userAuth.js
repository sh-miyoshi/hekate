import { AuthHandler } from '../plugins/auth.js'

export default async function(context) {
  const handler = new AuthHandler(context)

  const loginProject = window.localStorage.getItem('login_project')
  if (!loginProject) {
    handler.Login(context.params.name)
    return
  }

  const res = await handler.GetToken()
  if (!res.ok) {
    if (res.statusCode >= 400 && res.statusCode < 500) {
      context.redirect('/user/project/' + loginProject)
    } else {
      context.error({
        message: res.message,
        statusCode: 500
      })
    }
  }
}
