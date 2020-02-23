import { AuthHandler } from '../plugins/auth.js'

export default async function(context) {
  const expiresIn = window.localStorage.getItem('expires_in')
  const refreshToken = window.localStorage.getItem('refresh_token')
  if (!expiresIn || !refreshToken) {
    context.redirect('/')
    return
  }

  const now = Date.now() / 1000
  if (now >= expiresIn) {
    // access token was expired, so try to refresh
    const handler = new AuthHandler(context)
    const res = await handler.TokenRefresh(refreshToken)
    if (!res.ok) {
      if (res.statusCode >= 400 && res.statusCode < 500) {
        context.redirect('/')
      } else {
        context.error({
          message: res.message,
          statusCode: 500
        })
      }
    }
  }
}
