<template>
  <div>
    <p>redirecting ...</p>
  </div>
</template>

<script>
export default {
  middleware: 'auth',
  created() {
    this.$router.replace('/admin/home')
  }
}
</script>
