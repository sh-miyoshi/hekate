<template>
  <header class="c-header">
    <!-- logo start -->
    <nuxt-link to="/home" class="logo">HEKATE</nuxt-link>
    <!-- logo end -->

    <div class="">
      <ul class="">
        <li><button class="" @click="logout">Logout</button></li>
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
