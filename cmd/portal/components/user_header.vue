<template>
  <header class="c-header">
    <!-- logo start -->
    <img src="~/assets/img/logo.png" class="logo" />
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
    async logout() {
      const project = window.localStorage.getItem('login_project')
      await this.$auth.Logout()

      const to = '/user/project/' + project
      if (this.$router.currentRoute.path === to) {
        window.location.reload(true)
      } else {
        this.$router.push(to)
      }
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
