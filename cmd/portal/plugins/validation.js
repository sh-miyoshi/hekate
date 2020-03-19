function ValidateUserName(name) {
  if (typeof name !== 'string') {
    return { ok: false, msg: 'invalid name type.' }
  }

  if (name.length < 3 || name.length >= 64) {
    return { ok: false, msg: 'the length of name must be 3 to 63.' }
  }

  return { ok: true, msg: '' }
}

export default ({ app }, inject) => {
  inject('ValidateUserName', (string) => ValidateUserName(string))
}
