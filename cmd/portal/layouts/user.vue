<template>
  <div>
    <Sidebar />
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
import Header from '@/components/header.vue'
import Sidebar from '@/components/userSidebar.vue'

export default {
  components: {
    Header,
    Sidebar
  }
}
</script>

<style scoped>
.main {
  padding: 15px;
}
</style>
