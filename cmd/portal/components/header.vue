<template>
  <header class="c-header">
    <!-- logo start -->
    <nuxt-link to="/home" class="c-header-nav">
      <img src="~/assets/img/logo.png" class="logo" />
    </nuxt-link>
    <!-- logo end -->

    <div class="c-header-nav ml-auto mr-4">
      <button class="btn btn-primary px-4 btn-block" @click="logout">
        Logout
      </button>
    </div>
  </header>
</template>

<script>
export default {
  methods: {
    logout() {
      // TODO(logout and redirect to top)
    }
  }
}
</script>

<style scoped>
.logo {
  width: 120px;
  padding-left: 10px;
}
</style>
