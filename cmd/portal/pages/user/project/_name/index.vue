<template>
  <div>
    <p>redirecting ...</p>
  </div>
</template>

<script>
export default {
  layout: 'user_default',
  // TODO middleware: 'auth',
  mounted() {
    const project = window.localStorage.getItem('login_project')
    this.$router.replace('/user/project/' + project + '/home')
  }
}
</script>
