<template>
  <div class="container index-window">
    <div class="c-wrapper">
      <div class="c-body">
        <h4>Welcome to Hekate!</h4>
        <div class="card-body">
          <button class="btn btn-primary" @click="login">
            Login
          </button>
        </div>
      </div>
    </div>
    <footer class="c-footer"></footer>
  </div>
</template>

<script>
export default {
  layout: 'index',
  methods: {
    login() {
      this.$auth.Login(this.$route.params.name)
    }
  }
}
</script>
