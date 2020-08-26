import fs from 'fs'

export default {
  mode: 'spa',
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
  loading: { color: '#fff' },
  css: [
    '@/assets/css/bootstrap.min.css',
    '@/assets/css/coreui.min.css',
    '@/assets/css/style.css'
  ],
  plugins: ['~/plugins/auth.js', '~/plugins/api.js', '~/plugins/validation.js'],
  buildModules: [
    '@nuxtjs/eslint-module'
  ],
  modules: [
    '@nuxtjs/axios',
    '@nuxtjs/dotenv',
    '@nuxtjs/font-awesome',
    'bootstrap-vue/nuxt'
  ],
  axios: {},
  build: {},
  server: {
    host: '0.0.0.0',
    port: process.env.HEKATE_PORTAL_PORT,
    https: {
      key: fs.readFileSync('/hekatesecret/tls.key'),
      cert: fs.readFileSync('/hekatesecret/tls.crt')
    }
  },
  env: {
    HEKATE_SERVER_ADDR:
      process.env.HEKATE_SERVER_ADDR || 'http://localhost:18443',
    HEKATE_PORTAL_HOST: process.env.HEKATE_PORTAL_HOST || 'localhost',
    HEKATE_PORTAL_PORT: process.env.HEKATE_PORTAL_PORT || '3000',
    https: true
  }
}
