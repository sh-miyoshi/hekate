import { AuthHandler } from '../plugins/auth.js'

export default async function(context) {
  let project = ''
  let redirect = '/admin'
  if (context.route.path.includes('/user/project/')) {
    const values = context.route.path.split('/')
    if (values.length >= 4) {
      // ignore extra path
      project = values[3]
      redirect = '/user/project/' + project
    }
  }

  const handler = new AuthHandler(context)

  const loginProject = window.localStorage.getItem('login_project')
  if (loginProject == null || project !== loginProject) {
    handler.Login(project)
    return
  }

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
