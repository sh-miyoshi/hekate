export function ValidateUserName(name) {
  if (typeof name !== 'string') {
    return { ok: false, message: 'Invalid name type.' }
  }

  if (name.length < 3 || name.length >= 64) {
    return { ok: false, message: 'The length of name must be 3 to 63.' }
  }

  return { ok: true, message: '' }
}

export function ValidateClientID(id) {
  if (typeof id !== 'string') {
    return { ok: false, message: 'Invalid client id type.' }
  }

  const pattern = /^[a-z][a-z0-9\-._]{2,62}$/
  if (!id.match(pattern)) {
    return { ok: false, message: 'Invalid client id format.' }
  }
  return { ok: true, message: '' }
}

export function ValidateRoleName(name) {
  if (typeof name !== 'string') {
    return { ok: false, message: 'Invalid name type.' }
  }

  if (name.length < 3 || name.length >= 64) {
    return { ok: false, message: 'The length of name must be 3 to 63.' }
  }

  return { ok: true, message: '' }
}
