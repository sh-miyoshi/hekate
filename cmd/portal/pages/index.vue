<template>
  <div class="container index-window">
    <h1>JWT-Server</h1>
    <p><i>simple identity management server</i></p>
    <button class="btn btn-compose" @click="login">
      Login to admin console
    </button>
    {{ path }}
  </div>
</template>

<script>
export default {
  data() {
    return {
      path: `${process.env.OIDC_PATH}`
    }
  },
  methods: {
    login() {
      this.$auth.loginWith('jwtserver')
    }
  }
}
</script>
