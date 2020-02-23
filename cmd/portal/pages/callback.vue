<template>
  <div>
    <p>call back ...</p>
  </div>
</template>

<script>
export default {
  layout: 'index',
  created() {
    this.$auth.AuthCode(this.$route.query.code)
  }
}
</script>
