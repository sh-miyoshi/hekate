export class Config {
  constructor(configFilePath) {
    this.config = {
      SERVER_ADDR: 'http://localhost:18443',
      PORTAL_ADDR: 'http://localhost:3000',
      LOGIN_PROJECT: 'master',
      CLIENT_ID: 'portal'
    }

    // TODO set from configFile
  }

  get() {
    return this.config
  }
}

export default (context, inject) => {
  inject('config', new Config(process.env.CONFIG_FILE))
}
