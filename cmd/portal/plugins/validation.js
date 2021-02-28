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

export function ValidateUserCode(code) {
  if (typeof code !== 'string') {
    return { ok: false, message: 'The code is not string' }
  }

  if (code.length !== 6) {
    return { ok: false, message: 'The length of code is not 6' }
  }

  // all char is only digit
  for (let i = 0; i < code.length; i++) {
    if (code[i] < '0' || code[i] > '9') {
      return { ok: false, message: 'The code contains non-numeric characters' }
    }
  }

  return { ok: true, message: '' }
}
