import { AuthHandler } from '../plugins/auth.js'

export default async function(context) {
  const expiresIn = window.localStorage.getItem('expires_in')
  if (!expiresIn) {
    context.redirect('/')
    return
  }

  const handler = new AuthHandler(context)
  const res = await handler.GetToken()
  if (!res.ok) {
    if (res.statusCode >= 400 && res.statusCode < 500) {
      context.redirect('/admin')
    } else {
      context.error({
        message: res.message,
        statusCode: 500
      })
    }
  }
}
