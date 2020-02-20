<template>
  <div class="wrapper content">
    Please select item from menu bar
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  methods: {
    logout() {
      this.$auth.logout()
    }
  }
}
</script>
