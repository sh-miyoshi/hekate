<template>
  <div class="container index-window">
    <h1 class="h1">HEKATE</h1>
    <p><i>simple identity management server</i></p>
    <button class="btn btn-primary px-4 btn-block" @click="login">
      Login to admin console
    </button>
  </div>
</template>

<script>
export default {
  layout: 'index',
  methods: {
    login() {
      this.$auth.Login()
    }
  }
}
</script>
