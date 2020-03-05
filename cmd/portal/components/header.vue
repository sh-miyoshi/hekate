<template>
  <header class="header black-bg">
    <!-- logo start -->
    <nuxt-link to="/home" class="logo">hekate</nuxt-link>
    <!-- logo end -->

    <div class="top-menu">
      <ul class="nav pull-right top-menu">
        <li><button class="logout" @click="logout">Logout</button></li>
      </ul>
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
