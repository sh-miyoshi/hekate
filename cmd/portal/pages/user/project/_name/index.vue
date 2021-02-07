<template>
  <div>
    <h4>Welcome to Hekate!</h4>
    <button class="btn btn-primary" @click="login">
      Login
    </button>
  </div>
</template>

<script>
export default {
  layout: 'index',
  methods: {
    login() {
      this.$store.commit('setCurrentProject', this.$route.params.id)
      this.$auth.Login()
    }
  }
}
</script>
