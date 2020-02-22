<template>
  <div class="container index-window">
    <h1>JWT-Server</h1>
    <p><i>simple identity management server</i></p>
    <button class="btn btn-compose" @click="login">
      Login to admin console
    </button>
  </div>
</template>

<script>
export default {
  layout: 'index',
  methods: {
    login() {
      // TODO(login to system)
    }
  }
}
</script>
