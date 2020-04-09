<template>
  <div>
    <div class="c-wrapper">
      <Header />
      <div class="c-body main">
        <nuxt />
      </div>
      <footer class="c-footer"></footer>
    </div>
  </div>
</template>

<script>
import Header from '@/components/user_header.vue'

export default {
  components: {
    Header
  }
}
</script>

<style scoped>
.main {
  padding: 15px;
}
</style>
