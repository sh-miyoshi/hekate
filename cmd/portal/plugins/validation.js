function ValidateUserName(name) {
  if (typeof name !== 'string') {
    return { ok: false, msg: 'Invalid name type.' }
  }

  if (name.length < 3 || name.length >= 64) {
    return { ok: false, msg: 'The length of name must be 3 to 63.' }
  }

  return { ok: true, msg: '' }
}

function ValidateClientID(id) {
  if (typeof id !== 'string') {
    return { ok: false, msg: 'Invalid client id type.' }
  }

  const pattern = /^[a-z][a-z0-9\-._]{2,62}$/
  if (!id.match(pattern)) {
    return { ok: false, msg: 'Invalid client id format.' }
  }
  return { ok: true, msg: '' }
}

export default ({ app }, inject) => {
  inject('ValidateUserName', (string) => ValidateUserName(string))
  inject('ValidateClientID', (string) => ValidateClientID(string))
}
