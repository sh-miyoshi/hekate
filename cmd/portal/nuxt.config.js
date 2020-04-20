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
  css: [
    '@/assets/css/bootstrap.min.css',
    '@/assets/css/coreui.min.css',
    '@/assets/css/style.css'
  ],
  /*
   ** Plugins to load before mounting the App
   */
  plugins: ['~/plugins/auth.js', '~/plugins/api.js', '~/plugins/validation.js'],
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
    // Doc: https://axios.nuxtjs.org/usage
    '@nuxtjs/axios',
    // Doc: https://github.com/nuxt-community/dotenv-module
    '@nuxtjs/dotenv',
    '@nuxtjs/font-awesome',
    'bootstrap-vue/nuxt'
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
    port: process.env.HEKATE_PORTAL_PORT
  },

  env: {
    HEKATE_SERVER_ADDR:
      process.env.HEKATE_SERVER_ADDR || 'http://localhost:18443',
    HEKATE_PORTAL_HOST: process.env.HEKATE_PORTAL_HOST || 'localhost',
    HEKATE_PORTAL_PORT: process.env.HEKATE_PORTAL_PORT || '3000',
    // https: {}, // TODO(set params if run as https)
    SYSTEM_ROLES: [
      'read-cluster',
      'write-cluster',
      'read-project',
      'write-project',
      'read-role',
      'write-role',
      'read-user',
      'write-user',
      'read-client',
      'write-client',
      'read-customrole',
      'write-customrole'
    ]
  }
}
