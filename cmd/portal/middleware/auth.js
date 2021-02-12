import { AuthHandler } from '../plugins/auth.js'

export default async function(context) {
  const handler = new AuthHandler(context)

  const loginProject = window.localStorage.getItem('login_project')
  if (!loginProject) {
    handler.Login(process.env.LOGIN_PROJECT)
    return
  }

  const res = await handler.GetToken()
  if (!res.ok) {
    if (res.statusCode >= 400 && res.statusCode < 500) {
      context.redirect('/')
    } else {
      context.error({
        message: res.message,
        statusCode: 500
      })
    }
  } else {
    // Check role
    const roles = handler.GetUserSystemRoles()
    if (!roles.includes('read-cluster')) {
      context.error({
        message: 'user do not have permission: read-cluster',
        statusCode: 500
      })
    }
  }
}
