export default {
  mode: 'spa',
  /*
   ** Headers of the page
   */
  head: {
    title: process.env.npm_package_name || '',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      {
        hid: 'description',
        name: 'description',
        content: process.env.npm_package_description || ''
      }
    ],
    link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }]
  },
  /*
   ** Customize the progress-bar color
   */
  loading: { color: '#fff' },
  /*
   ** Global CSS
   */
  css: ['@/assets/css/bootstrap.min.css', '@/assets/css/style.css'],
  /*
   ** Plugins to load before mounting the App
   */
  plugins: [],
  /*
   ** Nuxt.js dev-modules
   */
  buildModules: [
    // Doc: https://github.com/nuxt-community/eslint-module
    '@nuxtjs/eslint-module'
  ],
  /*
   ** Nuxt.js modules
   */
  modules: [
    '@nuxtjs/auth',
    // Doc: https://axios.nuxtjs.org/usage
    '@nuxtjs/axios',
    // Doc: https://github.com/nuxt-community/dotenv-module
    '@nuxtjs/dotenv'
  ],
  /*
   ** Axios module configuration
   ** See https://axios.nuxtjs.org/options
   */
  axios: {},
  /*
   ** Build configuration
   */
  build: {
    /*
     ** You can extend webpack config here
     */
    extend(config, ctx) {
      // Run ESLint on save
      if (ctx.isDev && ctx.isClient) {
        config.module.rules.push({
          enforce: 'pre',
          test: /\.(js|vue)$/,
          loader: 'eslint-loader',
          exclude: /(node_modules)/,
          options: {
            fix: true
          }
        })
      }
    }
  },

  server: {
    host: '0.0.0.0',
    port: '3000'
  },

  env: {
    serverAddr: 'http://localhost:8080'
  },

  auth: {
    redirect: {
      login: '/',
      logout: '/',
      callback: '/callback',
      home: '/home',
    },
    strategies: {
      jwtserver: {
        _scheme: 'oauth2',
        authorization_endpoint: 'http://localhost:8080/api/v1/project/master/openid-connect/auth',
        userinfo_endpoint: 'http://localhost:8080/api/v1/project/master/openid-connect/userinfo',
        scope: ['openid'],
        access_token_endpoint: 'http://localhost:8080/api/v1/project/master/openid-connect/token',
        response_type: 'code',
        token_type: 'Bearer',
        redirect_uri: "http://localhost:3000/callback",
        client_id: 'admin-cli',
      },
      github: {
        client_id: "",
        client_secret: ""
      }
    },
  }
}
