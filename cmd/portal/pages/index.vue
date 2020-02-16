<template>
  <div class="container">
    <h2>Welcom to JWT-Server</h2>
    <button @click="login">Login</button>
    <button @click="login_github">Github Login</button>
  </div>
</template>

<script>
export default {
  methods: {
    login() {
      this.$auth.loginWith('jwtserver')
    },
    login_github() {
      this.$auth.loginWith('github')
    }
  }
}
</script>
